package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Log      Log
		Telegram Telegram
		Mongo    Mongo
		GPT      GPT
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"info"`
		Type  string `env:"LOG_TYPE" env-default:"text"`
	}

	Telegram struct {
		Token  string `env:"TELEGRAM_TOKEN" env-required:"true" `
	}

	Mongo struct {
		URL string `env:"MONGO_URL" env-required:"true"`
		DB  string `env:"MONGO_DB" env-required:"true"`
	}

	GPT struct {
		Model  string `env:"GPT_MODEL" env-required:"true"`
		Token  string `env:"GPT_TOKEN" env-required:"true"`
		Prompt string `env:"GPT_PROMPT" env-required:"true"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
