package on_voice

import (
	"context"
	"fmt"
	"io"
	"surgu-ai-chat-bot/internal/ai"

	"github.com/samber/lo"
	tele "gopkg.in/telebot.v4"
)

type Handler struct {
	bot      *tele.Bot
	speecher speecher
	ai       aiService
	webUrl   string
	storage  storage
}

func New(webUrl string, bot *tele.Bot, speecher speecher, ai aiService, storage storage) *Handler {
	return &Handler{
		bot:      bot,
		speecher: speecher,
		ai:       ai,
		webUrl:   webUrl,
		storage:  storage,
	}
}

func (h *Handler) Handle(c tele.Context) error {
	ctx := context.Background()

	voice := c.Message().Voice
	if voice == nil {
		return nil
	}

	if err := c.Send("⏳ Принял ваш запрос! Обрабатываю..."); err != nil {
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

	question, err := h.speecher.SpeechToText(ctx, buffer)
	if err != nil {
		return fmt.Errorf("h.speecher.SpeechToText: %w", err)
	}

	answer, err := h.ai.Answer(ctx, question)
	if err != nil {
		return fmt.Errorf("h.ai.Answer: %w", err)
	}

	if err := h.storage.LogVoiceQuestion(ctx, c.Chat().ID, question, lo.Map(answer.Files, func(f ai.File, _ int) string {
		return f.Name
	})); err != nil {
		return fmt.Errorf("h.storage.LogVoiceQuestion: %w", err)
	}

	answerText := fmt.Sprintf("✅ Готово! \nВаш вопрос: \"%s\"\n\nВот что мне удалось найти:\n\n%s\n\nНужна ещё помощь? Просто напишите новый запрос!", question, answer.Answer)
	if len(answer.Files) == 0 {
		return c.Send(answerText)
	}

	files, err := h.storage.GetFilesByNames(ctx, lo.Map(answer.Files, func(f ai.File, _ int) string {
		return f.Name
	}))
	if err != nil {
		return fmt.Errorf("h.storage.GetFilesByNames: %w", err)
	}

	markup := &tele.ReplyMarkup{}
	var rows []tele.Row
	for _, f := range files {
		rows = append(rows, markup.Row(markup.URL(f.ShortName, f.Url)))
	}
	markup.Inline(rows...)

	return c.Send(answerText, markup)
}
