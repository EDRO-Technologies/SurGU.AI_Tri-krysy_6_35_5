package on_text

import (
	"context"
	"fmt"
	"surgu-ai-chat-bot/internal/ai"

	"github.com/samber/lo"
	tele "gopkg.in/telebot.v4"
)

type Handler struct {
	bot     *tele.Bot
	ai      aiService
	storage storage
}

func New(bot *tele.Bot, ai aiService, storage storage) *Handler {
	return &Handler{
		bot:     bot,
		ai:      ai,
		storage: storage,
	}
}

func (h *Handler) Handle(c tele.Context) error {
	ctx := context.Background()

	if err := c.Send("⏳ Принял ваш запрос! Обрабатываю..."); err != nil {
		return fmt.Errorf("c.Send: %w", err)
	}

	question := c.Message().Text
	answer, err := h.ai.Answer(ctx, question)
	if err != nil {
		return fmt.Errorf("h.ai.Answer: %w", err)
	}

	if err := h.storage.LogTextQuestion(ctx, c.Chat().ID, question, lo.Map(answer.Files, func(f ai.File, _ int) string {
		return f.Name
	})); err != nil {
		return fmt.Errorf("h.storage.LogTextQuestion: %w", err)
	}

	answerText := fmt.Sprintf("✅ Готово! Вот что мне удалось найти:\n\n%s\n\nНужна ещё помощь? Просто отправь новый запрос!", answer.Answer)
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
