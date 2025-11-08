package on_text

import (
	"context"
	"surgu-ai-chat-bot/internal/ai"
)

type aiService interface {
	Answer(ctx context.Context, question string) (ai.Response, error)
}
