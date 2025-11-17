package logging

import (
    "io"
    "os"

    "github.com/lmittmann/tint"
    "log/slog"
)

// New creates a tint-backed slog.Logger according to environment.
func New(env string) *slog.Logger {
    opts := &tint.Options{
        Level: slog.LevelInfo,
    }

    if env == "development" {
        opts.NoColor = false
    } else {
        opts.NoColor = true
    }

    handler := tint.NewHandler(io.Discard, opts)
    handler = tint.NewHandler(os.Stdout, opts)
    return slog.New(handler)
}

