package config

import "github.com/caarlos0/env/v11"

type Config struct {
	OpenAIKey string `env:"OPENAI_API_KEY,required"`
	Model     string `env:"OPENAI_MODEL,required"`
	Port      string `env:"PORT,required"`
	DBPath    string `env:"DB_PATH,required"`
	BaseURL   string `env:"OPENAI_BASE_URL"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
