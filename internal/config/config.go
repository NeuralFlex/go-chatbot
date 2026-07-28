package config

import "github.com/caarlos0/env/v11"

type Config struct {
	OpenAIKey   string `env:"OPENAI_API_KEY,required"`
	Model       string `env:"OPENAI_MODEL,required"`
	Port        string `env:"PORT,required"`
	DatabaseURL string `env:"DATABASE_URL,required"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
