package check_privacy_policy

import (
	"context"
	"fmt"

	tele "gopkg.in/telebot.v4"
)

type Middleware struct {
	storage storage
}

func New(storage storage) *Middleware {
	return &Middleware{storage: storage}
}

func (m *Middleware) Func(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		signed, err := m.storage.PrivacyPolicySigned(context.Background(), c.Chat().ID)
		if err != nil {
			return fmt.Errorf("m.storage.PrivacyPolicySigned: %w", err)
		}

		if !signed {
			return c.Send("❌ Без согласия с политикой конфиденциальности я не могу обрабатывать ваши запросы.\n\nЕсли передумаете, используйте команду /start для повторного запуска бота и нажмите \"Прочитано\".")
		}

		return next(c)
	}
}
