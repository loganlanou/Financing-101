package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Config centralizes runtime configuration sourced from environment variables.
type Config struct {
	Env               string
	HTTPAddr          string
	PublicURL         string
	DatabasePath      string
	JWTAudience       string
	JWTIssuer         string
	ClerkSecretKey    string
	StripeSecretKey   string
	ShipStationAPIKey string
	ShipStationSecret string
	SendGridAPIKey    string
	SigningKey        string
	RequestTimeout    time.Duration
	NewsFeeds         []string
	NewsPollInterval  time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Env:               getEnv("APP_ENV", "development"),
		HTTPAddr:          getEnv("HTTP_ADDR", ":8080"),
		PublicURL:         getEnv("PUBLIC_URL", "http://localhost:8080"),
		DatabasePath:      getEnv("DATABASE_PATH", "data/app.db"),
		JWTAudience:       getEnv("JWT_AUDIENCE", "financing101"),
		JWTIssuer:         getEnv("JWT_ISSUER", "financing101"),
		ClerkSecretKey:    os.Getenv("CLERK_SECRET_KEY"),
		StripeSecretKey:   os.Getenv("STRIPE_SECRET_KEY"),
		ShipStationAPIKey: getEnv("SHIPSTATION_API_KEY", ""),
		ShipStationSecret: getEnv("SHIPSTATION_SECRET", ""),
		SendGridAPIKey:    getEnv("SENDGRID_API_KEY", ""),
		SigningKey:        getEnv("SIGNING_KEY", "insecure-local-key"),
	}

	timeoutStr := getEnv("REQUEST_TIMEOUT", "4s")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid REQUEST_TIMEOUT: %w", err)
	}

	cfg.RequestTimeout = timeout

	feeds := getEnv("NEWS_FEEDS", "https://finance.yahoo.com/news/rssindex,https://feeds.a.dj.com/rss/RSSMarketsMain")
	cfg.NewsFeeds = splitAndClean(feeds)

	pollStr := getEnv("NEWS_POLL_INTERVAL", "30m")
	pollDuration, err := time.ParseDuration(pollStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid NEWS_POLL_INTERVAL: %w", err)
	}
	cfg.NewsPollInterval = pollDuration

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}

func splitAndClean(raw string) []string {
	out := make([]string, 0)
	for _, part := range strings.Split(raw, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}
