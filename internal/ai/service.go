package ai

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
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

func (s *Service) Answer(_ context.Context, question string) (Response, error) {
	body, err := json.Marshal(Request{
		Question: question,
	})
	if err != nil {
		return Response{}, fmt.Errorf("json.Marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.baseUrl+"/ask", bytes.NewBuffer(body))
	if err != nil {
		return Response{}, fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("s.httpClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("status code not 200: %d", resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Response{}, fmt.Errorf("json.NewDecoder: Decode: %w", err)
	}

	return response, nil
}
