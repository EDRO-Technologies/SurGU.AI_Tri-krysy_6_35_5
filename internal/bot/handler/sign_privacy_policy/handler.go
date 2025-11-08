package sign_privacy_policy

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
	if err := h.storage.SignPrivacyPolicy(context.Background(), c.Chat().ID); err != nil {
		return fmt.Errorf("h.storage.UpdateUserPrivacyPolicySign: %w", err)
	}

	if err := c.Respond(); err != nil {
		return fmt.Errorf("h.Respond: %w", err)
	}

	return nil
}
