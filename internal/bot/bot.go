package bot

import (
	"fmt"
	"log"
	"strings"
	"surgu-ai-chat-bot/internal/ai"
	"surgu-ai-chat-bot/internal/bot/handler/on_start"
	"surgu-ai-chat-bot/internal/bot/handler/on_text"
	"surgu-ai-chat-bot/internal/bot/handler/on_voice"
	"surgu-ai-chat-bot/internal/bot/handler/sign_privacy_policy"
	"surgu-ai-chat-bot/internal/speacher"
	"surgu-ai-chat-bot/internal/storage"
	"time"

	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	bot       *tele.Bot
	storage   *storage.Storage
	speecher  *speacher.Service
	aiService *ai.Service
	appUrl    string
}

func MustNew(token, appUrl string, storage *storage.Storage, speecher *speacher.Service, aiService *ai.Service) *Bot {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 1 * time.Second},
		OnError: func(err error, c tele.Context) {
			log.Printf("telebot error: %v, context: %v", err, c)
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
	b.bot.Handle(tele.OnText, on_text.New(b.bot, b.aiService).Handle)
	b.bot.Handle(tele.OnVoice, on_voice.New(b.appUrl, b.bot, b.speecher, b.aiService).Handle)

	signPrivacyPolicyHandler := sign_privacy_policy.New(b.bot, b.storage)
	b.bot.Handle(tele.OnCallback, func(c tele.Context) error {
		switch strings.TrimSpace(c.Callback().Data) {
		case "sign-privacy-policy":
			return signPrivacyPolicyHandler.Handle(c)
		}
		return nil
	})

	b.bot.Start()
	return nil
}

func (b *Bot) Stop() {
	b.bot.Stop()
}
