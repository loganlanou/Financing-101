package payments

import (
    "context"

    "log/slog"
)

// StripeClient wires webhook verification + checkout creation later.
type StripeClient struct {
    secret string
    log    *slog.Logger
}

func NewStripeClient(secret string, log *slog.Logger) *StripeClient {
    return &StripeClient{secret: secret, log: log}
}

func (c *StripeClient) HealthCheck(ctx context.Context) error {
    if c.secret == "" {
        c.log.Warn("stripe disabled", slog.String("reason", "missing secret"))
        return nil
    }

    c.log.Info("stripe ready for payments")
    return nil
}

