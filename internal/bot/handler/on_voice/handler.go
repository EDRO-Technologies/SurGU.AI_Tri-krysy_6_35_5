package on_voice

import (
	"context"
	"fmt"
	tele "gopkg.in/telebot.v4"
	"io"
	"surgu-ai-chat-bot/internal/bot/handler/utils/get_files_markup"
)

type Handler struct {
	bot      *tele.Bot
	speecher speecher
	ai       aiService
	webUrl   string
}

func New(webUrl string, bot *tele.Bot, speecher speecher, ai aiService) *Handler {
	return &Handler{
		bot:      bot,
		speecher: speecher,
		ai:       ai,
		webUrl:   webUrl,
	}
}

func (h *Handler) Handle(c tele.Context) error {
	ctx := context.Background()

	voice := c.Message().Voice
	if voice == nil {
		return nil
	}

	if err := c.Send("Ваш вопрос получен. Формирую ответ..."); err != nil {
		return fmt.Errorf("c.Send: %w", err)
	}

	file, err := h.bot.FileByID(voice.FileID)
	if err != nil {
		return fmt.Errorf("h.bot.FileByID: %w", err)
	}

	reader, err := h.bot.File(&file)
	if err != nil {
		return fmt.Errorf("h.bot.File: %w", err)
	}
	defer reader.Close()

	buffer, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}

	text, err := h.speecher.SpeechToText(ctx, buffer)
	if err != nil {
		return fmt.Errorf("h.speecher.SpeechToText: %w", err)
	}

	answer, err := h.ai.Answer(ctx, text)
	if err != nil {
		return fmt.Errorf("h.ai.Answer: %w", err)
	}

	answerText := fmt.Sprintf("Вопрос: %s\n\n%s", text, answer.Answer)
	if len(answer.Files) == 0 {
		return c.Send(answerText)
	}

	return c.Send(answerText, get_files_markup.GetMarkup(answer.Files))
}
