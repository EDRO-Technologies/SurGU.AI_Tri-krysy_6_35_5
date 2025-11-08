package bot

import (
	"context"
	"fmt"
	"strings"
	"surgu-ai-chat-bot/internal/ai"
	"surgu-ai-chat-bot/internal/bot/handler/on_start"
	"surgu-ai-chat-bot/internal/bot/handler/on_statistics"
	"surgu-ai-chat-bot/internal/bot/handler/on_text"
	"surgu-ai-chat-bot/internal/bot/handler/on_voice"
	"surgu-ai-chat-bot/internal/bot/handler/sign_privacy_policy"
	"surgu-ai-chat-bot/internal/bot/middleware/check_privacy_policy"
	"surgu-ai-chat-bot/internal/speacher"
	"surgu-ai-chat-bot/internal/storage"
	"surgu-ai-chat-bot/pkg/logger"
	"time"

	"gopkg.in/telebot.v4/middleware"

	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	bot       *tele.Bot
	storage   *storage.Storage
	speecher  *speacher.Service
	aiService *ai.Service
	appUrl    string
}

func MustNew(log logger.Log, token, appUrl string, storage *storage.Storage, speecher *speacher.Service, aiService *ai.Service) *Bot {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 1 * time.Second},
		OnError: func(err error, c tele.Context) {
			log.WithContext(context.Background()).WithError(err).WithFields(map[string]any{
				"user_id": c.Chat().ID,
			}).Error("telebot error")
			_ = c.Send("😔 Извините, не удалось обработать ваш запрос. Попробуйте немного позже.")
		},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		panic(fmt.Errorf("telebot.NewBot: %w", err))
	}

	return &Bot{
		appUrl:    appUrl,
		bot:       bot,
		storage:   storage,
		speecher:  speecher,
		aiService: aiService,
	}
}

func (b *Bot) Start() error {
	b.bot.Handle("/start", on_start.New(b.bot, b.storage).Handle)

	signPrivacyPolicyHandler := sign_privacy_policy.New(b.bot, b.storage)
	b.bot.Handle(tele.OnCallback, func(c tele.Context) error {
		switch strings.TrimSpace(c.Callback().Data) {
		case "sign-privacy-policy":
			return signPrivacyPolicyHandler.Handle(c)
		}
		return nil
	})

	adminOnly := b.bot.Group()
	adminOnly.Use(middleware.Whitelist(1069351042))
	adminOnly.Handle("/statistics", on_statistics.New(b.bot, b.storage).Handle)

	needSignedPrivacyPolicy := b.bot.Group()
	needSignedPrivacyPolicy.Use(check_privacy_policy.New(b.storage).Func)
	needSignedPrivacyPolicy.Handle(tele.OnText, on_text.New(b.bot, b.aiService, b.storage).Handle)
	needSignedPrivacyPolicy.Handle(tele.OnVoice, on_voice.New(b.appUrl, b.bot, b.speecher, b.aiService, b.storage).Handle)

	b.bot.Start()
	return nil
}

func (b *Bot) Stop() {
	b.bot.Stop()
}
