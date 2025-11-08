package on_voice

import (
	"context"
	"surgu-ai-chat-bot/internal/ai"
)

type speecher interface {
	SpeechToText(ctx context.Context, audioBuffer []byte) (string, error)
}

type aiService interface {
	Answer(ctx context.Context, question string) (ai.Response, error)
}
