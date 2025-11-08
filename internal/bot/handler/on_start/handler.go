package on_start

import (
	"context"
	"fmt"

	tele "gopkg.in/telebot.v4"
)

type Handler struct {
	bot     *tele.Bot
	storage storage
}

func New(bot *tele.Bot, storage storage) *Handler {
	return &Handler{
		bot:     bot,
		storage: storage,
	}
}

func (h *Handler) Handle(c tele.Context) error {
	ctx := context.Background()

	markup := &tele.ReplyMarkup{}
	markup.Inline(
		markup.Row(markup.URL("Политика конфиденциальности", "https://www.gazprom.ru/about/legal/policy-personal-data/")),
		markup.Row(markup.Data("Прочитано", "sign-privacy-policy")),
	)

	if err := h.storage.CreateUser(ctx, c.Chat().ID); err != nil {
		return fmt.Errorf("h.storage.CreateUser: %w", err)
	}

	return c.Send(
		"👋 Добро пожаловать!\n\nЯ бот по охране труда, который проводит обучение и ответит на вопросы по охране труда.\n\nПеред началом работы, пожалуйста, ознакомьтесь с политикой конфиденциальности.",
		markup,
		tele.ModeHTML,
	)
}
