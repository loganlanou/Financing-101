package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/loganlanou/Financing-101/db/migrations"
	"github.com/loganlanou/Financing-101/internal/auth"
	"github.com/loganlanou/Financing-101/internal/config"
	"github.com/loganlanou/Financing-101/internal/data"
	"github.com/loganlanou/Financing-101/internal/database"
	"github.com/loganlanou/Financing-101/internal/handlers"
	"github.com/loganlanou/Financing-101/internal/ingest"
	"github.com/loganlanou/Financing-101/internal/logging"
	"github.com/loganlanou/Financing-101/internal/mail"
	"github.com/loganlanou/Financing-101/internal/payments"
	"github.com/loganlanou/Financing-101/internal/server"
	"github.com/loganlanou/Financing-101/internal/services"
	"github.com/loganlanou/Financing-101/internal/shipping"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logging.New(cfg.Env)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, cfg, log); err != nil {
		log.Error("server exited", slog.Any("err", err))
		os.Exit(1)
	}
}

func run(ctx context.Context, cfg config.Config, log *slog.Logger) error {
	if err := os.MkdirAll(filepath.Dir(cfg.DatabasePath), 0o755); err != nil {
		return err
	}

	db, err := sql.Open("sqlite", cfg.DatabasePath)
	if err != nil {
		return err
	}
	defer db.Close()

	db.SetConnMaxLifetime(1 * time.Hour)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	goose.SetBaseFS(migrations.Files)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}
	if err := goose.Up(db, "."); err != nil {
		return err
	}

	queries := database.New(db)

	// Auxiliary integrations log health but remain optional.
	clerkClient := auth.NewClerkClient(cfg.ClerkSecretKey, cfg.PublicURL, log)
	stripeClient := payments.NewStripeClient(cfg.StripeSecretKey, log)
	shipClient := shipping.NewShipStationClient(cfg.ShipStationAPIKey, cfg.ShipStationSecret, log)
	mailClient := mail.NewSendGridClient(cfg.SendGridAPIKey, log)

	_ = clerkClient.HealthCheck(ctx)
	_ = stripeClient.HealthCheck(ctx)
	_ = shipClient.HealthCheck(ctx)
	_ = mailClient.HealthCheck(ctx)

	// Seed local data when the database is empty.
	if err := data.EnsureSeedData(ctx, queries, log, uuid.NewString); err != nil {
		return err
	}

	newsService := services.NewNewsService(log, queries)
	stockService := services.NewStockService(log, queries)
	tradeService := services.NewTradeService(log, queries)
	recService := services.NewRecommendationService(log, queries)
	learnService := services.NewLearnService(log, queries)

	newsIngestor := ingest.NewNewsIngestor(log, queries, cfg.NewsFeeds)
	if err := newsIngestor.Refresh(ctx, 20); err != nil {
		log.Warn("initial news ingest failed", slog.Any("err", err))
	}

	go func() {
		ticker := time.NewTicker(cfg.NewsPollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				refreshCtx, cancel := context.WithTimeout(context.Background(), cfg.RequestTimeout)
				if err := newsIngestor.Refresh(refreshCtx, 20); err != nil {
					log.Warn("scheduled news ingest failed", slog.Any("err", err))
				}
				cancel()
			}
		}
	}()

	srv := server.New(cfg, log)

	pagesHandler := handlers.NewPagesHandler(log, newsService, stockService, tradeService, recService)
	pagesHandler.RegisterRoutes(srv.Echo())

	return srv.Start(ctx)
}
