package on_text

import (
	"context"
	"fmt"
	tele "gopkg.in/telebot.v4"
	"surgu-ai-chat-bot/internal/bot/handler/utils/get_files_markup"
)

type Handler struct {
	bot *tele.Bot
	ai  aiService
}

func New(bot *tele.Bot, ai aiService) *Handler {
	return &Handler{
		bot: bot,
		ai:  ai,
	}
}

func (h *Handler) Handle(c tele.Context) error {
	ctx := context.Background()

	if err := c.Send("⏳ Принял ваш запрос! Обрабатываю..."); err != nil {
		return fmt.Errorf("c.Send: %w", err)
	}

	answer, err := h.ai.Answer(ctx, c.Message().Text)
	if err != nil {
		return fmt.Errorf("h.ai.Answer: %w", err)
	}

	answerText := fmt.Sprintf("✅ Готово! Вот что мне удалось найти:\n\n%s\n\nНужна ещё помощь? Просто напишите новый запрос!", answer.Answer)
	if len(answer.Files) == 0 {
		return c.Send(answerText)
	}

	return c.Send(answerText, get_files_markup.GetMarkup(answer.Files))
}
