package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"surgu-ai-chat-bot/internal/ai"
	"surgu-ai-chat-bot/internal/bot"
	"surgu-ai-chat-bot/internal/config"
	"surgu-ai-chat-bot/internal/speacher"
	"surgu-ai-chat-bot/internal/storage"
	"surgu-ai-chat-bot/pkg/logger"
	"surgu-ai-chat-bot/pkg/migrate"
	"surgu-ai-chat-bot/pkg/postgres"
	"time"
)

func main() {
	if err := realMain(); err != nil {
		fmt.Println(err)
	}
}

func realMain() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := config.MustNewConfig()

	log := logger.NewLogger(slog.LevelDebug, cfg.Env, os.Stdout)
	httpClient := &http.Client{Timeout: time.Second * 600}

	b := bot.MustNew(
		cfg.TgBotToken,
		cfg.WebAppUrl,
		storage.New(postgres.MustNew(cfg.DSN)),
		speacher.New("http://10.76.16.133:5000", httpClient),
		ai.New("http://10.76.16.133:5000", httpClient),
	)

	if err := migrate.Migrate(cfg.DSN); err != nil {
		return fmt.Errorf("migrate.Migrate: %w", err)
	}

	errorGroup, _ := errgroup.WithContext(ctx)

	errorGroup.Go(func() (err error) {
		<-ctx.Done()

		log.WithContext(ctx).Info("shutting down bot")
		b.Stop()

		return nil
	})

	errorGroup.Go(func() (err error) {
		log.WithContext(ctx).Info("starting bot")

		if err := b.Start(); err != nil {
			return fmt.Errorf("b.Start(): %w", err)
		}
		return nil
	})

	if err := errorGroup.Wait(); err != nil {
		return fmt.Errorf("errorGroup.Wait(): %w", err)
	}

	return nil
}
