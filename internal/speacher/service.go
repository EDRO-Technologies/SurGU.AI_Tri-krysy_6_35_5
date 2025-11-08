package speacher

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"mime/multipart"
	"net/http"
	ogg_to_waw "surgu-ai-chat-bot/pkg/ogg-to-waw"
)

type Service struct {
	baseUrl    string
	httpClient httpClient
}

func New(baseUrl string, httpClient httpClient) *Service {
	return &Service{
		baseUrl:    baseUrl,
		httpClient: httpClient,
	}
}

func (s *Service) SpeechToText(_ context.Context, buffer []byte) (string, error) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	part, err := writer.CreateFormFile("audio", "audio.wav")
	if err != nil {
		return "", fmt.Errorf("writer.CreateFormFile: %w", err)
	}

	wawBuffer, err := ogg_to_waw.Convert(buffer)
	if err != nil {
		return "", fmt.Errorf("ogg_to_waw.Convert: %w", err)
	}

	if _, err = part.Write(wawBuffer); err != nil {
		return "", fmt.Errorf("part.Write: %w", err)
	}

	if err = writer.Close(); err != nil {
		return "", fmt.Errorf("writer.Close: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.baseUrl+"/speech-to-text", &b)
	if err != nil {
		return "", fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("s.httpClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code not 200: %d", resp.StatusCode)
	}

	var response SpeechResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("json.NewDecoder: Decode: %w", err)
	}

	return response.Text, nil
}
