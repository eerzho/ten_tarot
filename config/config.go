package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Log
		HTTP
		Telegram
		Mongo
		Google
		GPT
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"info"`
	}

	HTTP struct {
		Port   string `env:"HTTP_PORT" env-default:"8080"`
		Domain string `env:"HTTP_DOMAIN" env-required:"true"`
	}

	Telegram struct {
		Port   string `env:"TELEGRAM_PORT" env-default:"8081"`
		Domain string `env:"TELEGRAM_DOMAIN" env-required:"true"`
		Token  string `env:"TELEGRAM_TOKEN" env-required:"true" `
	}

	Mongo struct {
		URL string `env:"MONGO_URL" env-required:"true"`
		DB  string `env:"MONGO_DB" env-default:"event_manager"`
	}

	Google struct {
		CalendarURL string `env:"GOOGLE_CALENDAR_URL" env-required:"true"`
	}

	GPT struct {
		Token  string `env:"GPT_TOKEN" env-required:"true"`
		Prompt string `env:"GPT_PROMPT" env-required:"true"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("./config::New: %w", err)
	}

	return cfg, nil
}
