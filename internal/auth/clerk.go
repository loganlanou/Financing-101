package auth

import (
    "context"

    "log/slog"
)

// ClerkClient is a lightweight placeholder for downstream auth workflows.
type ClerkClient struct {
    apiKey   string
    publicURL string
    log      *slog.Logger
}

func NewClerkClient(apiKey, publicURL string, log *slog.Logger) *ClerkClient {
    return &ClerkClient{apiKey: apiKey, publicURL: publicURL, log: log}
}

func (c *ClerkClient) HealthCheck(ctx context.Context) error {
    if c.apiKey == "" {
        c.log.Warn("clerk disabled", slog.String("reason", "missing api key"))
        return nil
    }

    c.log.Info("clerk ready", slog.String("public_url", c.publicURL))
    return nil
}

