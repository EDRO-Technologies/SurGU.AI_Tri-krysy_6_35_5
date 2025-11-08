package on_statistics

import (
	"context"
	storage_stats "surgu-ai-chat-bot/internal/storage"
)

type storage interface {
	GetStatistics(ctx context.Context) (storage_stats.Statistics, error)
}
