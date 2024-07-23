package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Log
		Crypter
		HTTP
		Telegram
		Mongo
		GPT
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"info"`
		Type  string `env:"LOG_TYPE" env-default:"text"`
	}

	Crypter struct {
		Key string `env:"CRYPTER_KEY" env-required:"true"`
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
		DB  string `env:"MONGO_DB" env-required:"true"`
	}

	GPT struct {
		Model  string `env:"GPT_MODEL" env-default:"gpt-3.5-turbo"`
		Token  string `env:"GPT_TOKEN" env-required:"true"`
		Prompt string `env:"GPT_PROMPT" env-required:"true"`
	}
)

func New() (*Config, error) {
	const op = "./config::New"

	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cfg, nil
}
