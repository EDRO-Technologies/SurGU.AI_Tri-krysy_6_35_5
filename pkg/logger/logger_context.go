package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
)

type LogCtx struct {
	log slogLogger
	ctx context.Context
}

func (l LogCtx) WithError(err error) LogCtx {
	if err == nil {
		return l
	}

	return LogCtx{
		log: l.log.With(slog.Any(errKey, err)),
		ctx: l.ctx,
	}
}

func (l LogCtx) WithFields(fields map[string]any) LogCtx {
	var attrs []any

	for k, f := range fields {
		attrs = append(attrs, slog.Any(k, f))
	}

	return LogCtx{
		log: l.log.With(attrs...),
		ctx: l.ctx,
	}
}

func (l LogCtx) Debug(msg string) {
	if !l.log.Enabled(l.ctx, slog.LevelDebug) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	_ = l.log.Handler().Handle(l.ctx, slog.NewRecord(time.Now(), slog.LevelDebug, msg, pcs[0]))
}

func (l LogCtx) Info(msg string) {
	if !l.log.Enabled(l.ctx, slog.LevelInfo) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	_ = l.log.Handler().Handle(l.ctx, slog.NewRecord(time.Now(), slog.LevelInfo, msg, pcs[0]))
}

func (l LogCtx) Warning(msg string) {
	if !l.log.Enabled(l.ctx, slog.LevelWarn) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	_ = l.log.Handler().Handle(l.ctx, slog.NewRecord(time.Now(), slog.LevelWarn, msg, pcs[0]))
}

func (l LogCtx) Error(msg string) {
	if !l.log.Enabled(l.ctx, slog.LevelError) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	_ = l.log.Handler().Handle(l.ctx, slog.NewRecord(time.Now(), slog.LevelError, msg, pcs[0]))
}

func (l LogCtx) Infof(template string, args ...any) {
	if !l.log.Enabled(l.ctx, slog.LevelInfo) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	_ = l.log.Handler().Handle(l.ctx, slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf(template, args...), pcs[0]))
}

func (l LogCtx) Errorf(template string, args ...any) {
	if !l.log.Enabled(l.ctx, slog.LevelError) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	_ = l.log.Handler().Handle(l.ctx, slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf(template, args...), pcs[0]))
}

func (l LogCtx) Fatalf(template string, args ...any) {
	if !l.log.Enabled(l.ctx, slog.LevelError) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	_ = l.log.Handler().Handle(l.ctx, slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf(template, args...), pcs[0]))
	os.Exit(1)
}
