package config

import (
	"errors"
	"os"
	"time"
)

type Config struct {
	Provider      string // openrouter | deepseek
	BaseURL       string
	Model         string
	APIKey        string
	DB_DSN        string
	MigrationsDir string // path to goose migrations
	HTTPAddr      string // ":8080"
	PromptPath    string // path to file with system prompt
	SystemPrompt  string // loaded prompt content
	AppReferer    string // optional: for OpenRouter leaderboards
	AppTitle      string // optional: for OpenRouter leaderboards
}

func Load() (Config, error) {
	cfg := Config{}
	cfg.Provider = getenv("PROVIDER", "openrouter")
	cfg.BaseURL = getenv("BASE_URL", defaultBase(cfg.Provider))
	cfg.Model = getenv("MODEL", defaultModel(cfg.Provider))
	cfg.APIKey = getenv("API_KEY", "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")
	cfg.DB_DSN = getenv("DB_DSN", "postgres://moleben:moleben@localhost:5432/moleben?sslmode=disable")
	cfg.MigrationsDir = getenv("MIGRATIONS_DIR", "./migrations")
	cfg.HTTPAddr = getenv("HTTP_ADDR", ":8080")
	cfg.PromptPath = getenv("PROMPT_PATH", "prompt")
	if cfg.APIKey == "" {
		return cfg, errors.New("API_KEY is required")
	}
	b, err := os.ReadFile(cfg.PromptPath)
	if err != nil {
		return cfg, err
	}
	cfg.SystemPrompt = string(b)
	cfg.AppReferer = getenv("APP_REFERER", "http://localhost:3000")
	cfg.AppTitle = getenv("APP_TITLE", "Moleben")
	return cfg, nil
}

func defaultBase(provider string) string {
	if provider == "openrouter" {
		return "https://openrouter.ai/api/v1"
	}
	return "https://api.deepseek.com"
}

func defaultModel(provider string) string {
	if provider == "openrouter" {
		return "nvidia/nemotron-nano-9b-v2:free"
	}
	return "deepseek-chat"
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// For future: central timeouts
func RequestTimeout() time.Duration { return 60 * time.Second }
