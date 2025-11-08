package on_text

import (
	"context"
	"surgu-ai-chat-bot/internal/ai"
	_storage "surgu-ai-chat-bot/internal/storage"
)

type aiService interface {
	Answer(ctx context.Context, question string) (ai.Response, error)
}

type storage interface {
	LogTextQuestion(ctx context.Context, userId int64, question string, fileNames []string) error
	GetFilesByNames(ctx context.Context, names []string) ([]_storage.File, error)
}
