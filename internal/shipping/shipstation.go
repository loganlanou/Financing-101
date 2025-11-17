package shipping

import (
    "context"

    "log/slog"
)

// ShipStationClient keeps shipping orchestration configuration.
type ShipStationClient struct {
    apiKey string
    secret string
    log    *slog.Logger
}

func NewShipStationClient(apiKey, secret string, log *slog.Logger) *ShipStationClient {
    return &ShipStationClient{apiKey: apiKey, secret: secret, log: log}
}

func (c *ShipStationClient) HealthCheck(ctx context.Context) error {
    if c.apiKey == "" || c.secret == "" {
        c.log.Warn("shipstation disabled", slog.String("reason", "missing credentials"))
        return nil
    }

    c.log.Info("shipstation ready")
    return nil
}

