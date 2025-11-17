package server

import (
    "context"
    "net/http"
    "time"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/loganlanou/Financing-101/internal/config"
    "log/slog"
)

// Server wraps Echo and lifecycle hooks.
type Server struct {
    app *echo.Echo
    cfg config.Config
    log *slog.Logger
}

func New(cfg config.Config, log *slog.Logger) *Server {
    e := echo.New()
    e.HideBanner = true
    e.HidePort = true

    e.Use(middleware.Recover())
    e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
        LogStatus: true,
        LogURI:    true,
        LogLatency: true,
        LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
            log.Info("http",
                slog.String("method", c.Request().Method),
                slog.String("path", v.URI),
                slog.Int("status", v.Status),
                slog.Duration("latency", v.Latency),
            )
            return nil
        },
    }))

    e.Static("/static", "web/static")

    return &Server{app: e, cfg: cfg, log: log}
}

func (s *Server) Echo() *echo.Echo {
    return s.app
}

func (s *Server) Start(ctx context.Context) error {
    srvErr := make(chan error, 1)

    go func() {
        if err := s.app.Start(s.cfg.HTTPAddr); err != nil && err != http.ErrServerClosed {
            srvErr <- err
        }
    }()

    select {
    case <-ctx.Done():
        shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        if err := s.app.Shutdown(shutdownCtx); err != nil {
            return err
        }
        return nil
    case err := <-srvErr:
        return err
    }
}

