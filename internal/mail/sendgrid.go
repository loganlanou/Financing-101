package mail

import (
    "context"

    "log/slog"
)

// SendGridClient stubs transactional email flows used for onboarding + alerts.
type SendGridClient struct {
    apiKey string
    log    *slog.Logger
}

func NewSendGridClient(apiKey string, log *slog.Logger) *SendGridClient {
    return &SendGridClient{apiKey: apiKey, log: log}
}

func (c *SendGridClient) HealthCheck(ctx context.Context) error {
    if c.apiKey == "" {
        c.log.Warn("sendgrid disabled", slog.String("reason", "missing api key"))
        return nil
    }

    c.log.Info("sendgrid ready for transactional email")
    return nil
}

