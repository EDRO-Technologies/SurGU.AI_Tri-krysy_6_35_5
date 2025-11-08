package on_statistics

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgtype"
	"github.com/samber/lo"
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

	stats, err := h.storage.GetStatistics(ctx)
	if err != nil {
		return fmt.Errorf("h.storage.GetStatistics: %w", err)
	}

	return c.Send(fmt.Sprintf(
		"Всего пользователей: %d\nВсего вопросов: %d (%d голосовых, %d текстом)\n\nПоследние 10 запросов:\n%s",
		stats.TotalUsers,
		stats.TotalTextQuestions+stats.TotalVoiceQuestions,
		stats.TotalVoiceQuestions,
		stats.TotalTextQuestions,
		strings.Join(lo.Map(stats.Last10Questions.Elements, func(q pgtype.Text, i int) string {
			return fmt.Sprintf("%d. \"%s\"", i+1, q.String)
		}), "\n"),
	))
}
