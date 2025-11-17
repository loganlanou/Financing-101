package main

import (
    "context"
    "database/sql"
    "log"
    "log/slog"
    "os"

    "github.com/google/uuid"
    "github.com/loganlanou/Financing-101/db/migrations"
    "github.com/loganlanou/Financing-101/internal/config"
    "github.com/loganlanou/Financing-101/internal/data"
    "github.com/loganlanou/Financing-101/internal/database"
    "github.com/pressly/goose/v3"
    _ "modernc.org/sqlite"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    db, err := sql.Open("sqlite", cfg.DatabasePath)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    goose.SetBaseFS(migrations.Files)
    if err := goose.SetDialect("sqlite3"); err != nil {
        log.Fatal(err)
    }
    if err := goose.Up(db, "db/migrations"); err != nil {
        log.Fatal(err)
    }

    queries := database.New(db)
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    if err := data.EnsureSeedData(context.Background(), queries, logger, uuid.NewString); err != nil {
        log.Fatal(err)
    }

    log.Println("Seed data ensured")
}

