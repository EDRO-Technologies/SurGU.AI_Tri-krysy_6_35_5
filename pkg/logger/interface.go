package logger

import "context"

type Log interface {
	WithContext(ctx context.Context) LogCtx
}
