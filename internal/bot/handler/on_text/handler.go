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

	if err := c.Send("Ваш вопрос получен. Формирую ответ..."); err != nil {
		return fmt.Errorf("c.Send: %w", err)
	}

	answer, err := h.ai.Answer(ctx, c.Message().Text)
	if err != nil {
		return fmt.Errorf("h.ai.Answer: %w", err)
	}

	if len(answer.Files) == 0 {
		return c.Send(answer.Answer)
	}

	return c.Send(answer.Answer, get_files_markup.GetMarkup(answer.Files))
}
