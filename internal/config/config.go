package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env        string
	ServerPort string
	TgBotToken string
	WebAppUrl  string
	StaticDir  string
	DSN        string
}

func MustNewConfig() Config {
	_ = godotenv.Load()

	return Config{
		Env:        mustGetEnv("ENVIRONMENT"),
		ServerPort: mustGetEnv("HTTP_PORT"),
		TgBotToken: mustGetEnv("TG_BOT_TOKEN"),
		WebAppUrl:  mustGetEnv("WEB_APP_URL"),
		StaticDir:  mustGetEnv("STATIC_DIR"), // TODO: заменить яндекс диск на local storage
		DSN:        mustGetEnv("DATABASE_URL"),
	}
}

func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("environment variable %s not set", key))
	}
	return value
}
