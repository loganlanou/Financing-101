package data

import (
    "context"
    "database/sql"
    "time"

    "github.com/loganlanou/Financing-101/internal/database"
    "log/slog"
)

// EnsureSeedData hydrates the local sqlite db with demo insights for SSR rendering.
func EnsureSeedData(ctx context.Context, queries *database.Queries, log *slog.Logger, uuidFn func() string) error {
    if err := seedNews(ctx, queries, log, uuidFn); err != nil {
        return err
    }
    if err := seedStocks(ctx, queries, log, uuidFn); err != nil {
        return err
    }
    if err := seedTrades(ctx, queries, log, uuidFn); err != nil {
        return err
    }
    if err := seedRecommendations(ctx, queries, log, uuidFn); err != nil {
        return err
    }
    return nil
}

func seedNews(ctx context.Context, queries *database.Queries, log *slog.Logger, uuidFn func() string) error {
    rows, err := queries.ListLatestNews(ctx, 1)
    if err == nil && len(rows) > 0 {
        return nil
    }

    samples := []database.InsertNewsArticleParams{
        {
            ID:             uuidFn(),
            Title:          "AI Chips Surge as Earnings Beat Forecasts",
            Source:         "CNBC",
            Summary:        "NVIDIA extends gains after highlighting hyperscaler demand and AI enterprise adoption.",
            SentimentScore: 0.78,
            Trend:          "bullish",
            Tickers:        "NVDA, MSFT",
            Url:            "https://www.cnbc.com/",
            PublishedAt:    time.Now().Add(-2 * time.Hour),
        },
        {
            ID:             uuidFn(),
            Title:          "Fed Minutes show soft landing hopes",
            Source:         "WSJ",
            Summary:        "Officials continue to weigh one final hike while acknowledging cooling inflation prints.",
            SentimentScore: 0.21,
            Trend:          "neutral",
            Tickers:        "SPY, QQQ",
            Url:            "https://www.wsj.com/",
            PublishedAt:    time.Now().Add(-6 * time.Hour),
        },
        {
            ID:             uuidFn(),
            Title:          "Energy transition update from COP summit",
            Source:         "Bloomberg",
            Summary:        "Policy makers outline incentives for clean grid build-outs as oil majors pledge new CAPEX.",
            SentimentScore: -0.12,
            Trend:          "bearish",
            Tickers:        "XOM, ENPH",
            Url:            "https://www.bloomberg.com/",
            PublishedAt:    time.Now().Add(-12 * time.Hour),
        },
    }

    for _, sample := range samples {
        if err := queries.InsertNewsArticle(ctx, sample); err != nil {
            return err
        }
    }

    log.Info("seeded news dataset")
    return nil
}

func seedStocks(ctx context.Context, queries *database.Queries, log *slog.Logger, uuidFn func() string) error {
    rows, err := queries.ListStockSnapshots(ctx, 1)
    if err == nil && len(rows) > 0 {
        return nil
    }

    now := time.Now()
    payloads := []database.InsertStockSnapshotParams{
        {
            ID:            uuidFn(),
            Symbol:        "NVDA",
            Name:          "NVIDIA Corp",
            Sector:        sqlNullString("Technology"),
            Industry:      sqlNullString("Semiconductors"),
            Change30:      14.2,
            Change90:      32.8,
            Change365:     155.4,
            VsSp50030:     8.4,
            VsSp50090:     21.3,
            VsSp500365:    124.9,
            Conviction:    "High",
            Thesis:        "Accelerated compute demand from GenAI pipelines.",
            UpdatedAt:     now,
        },
        {
            ID:            uuidFn(),
            Symbol:        "MSFT",
            Name:          "Microsoft Corp",
            Sector:        sqlNullString("Technology"),
            Industry:      sqlNullString("Software"),
            Change30:      6.5,
            Change90:      18.1,
            Change365:     38.0,
            VsSp50030:     2.0,
            VsSp50090:     6.2,
            VsSp500365:    12.1,
            Conviction:    "Medium",
            Thesis:        "Copilot monetization and Azure AI workloads drive durable growth.",
            UpdatedAt:     now,
        },
        {
            ID:            uuidFn(),
            Symbol:        "TSLA",
            Name:          "Tesla Inc",
            Sector:        sqlNullString("Consumer Discretionary"),
            Industry:      sqlNullString("Automotive"),
            Change30:      -4.3,
            Change90:      5.1,
            Change365:     12.4,
            VsSp50030:     -7.1,
            VsSp50090:     -6.8,
            VsSp500365:    -13.5,
            Conviction:    "Watching",
            Thesis:        "Margin compression from price cuts offset by energy storage upside.",
            UpdatedAt:     now,
        },
    }

    for _, payload := range payloads {
        if err := queries.InsertStockSnapshot(ctx, payload); err != nil {
            return err
        }
    }

    log.Info("seeded stock snapshots")
    return nil
}

func seedTrades(ctx context.Context, queries *database.Queries, log *slog.Logger, uuidFn func() string) error {
    rows, err := queries.ListRecentTrades(ctx, 1)
    if err == nil && len(rows) > 0 {
        return nil
    }

    events := []database.InsertCongressTradeParams{
        {
            ID:             uuidFn(),
            Member:         "Nancy Pelosi",
            Party:          sqlNullString("D"),
            Chamber:        sqlNullString("House"),
            Symbol:         "NVDA",
            Action:         "Buy",
            Amount:         sqlNullString("$1M - $5M"),
            ExecutedAt:     time.Now().AddDate(0, 0, -7),
            DisclosureDate: time.Now().AddDate(0, 0, -2),
            Sentiment:      0.75,
            SourceUrl:      sqlNullString("https://senatestockwatcher.com"),
        },
        {
            ID:             uuidFn(),
            Member:         "Dan Crenshaw",
            Party:          sqlNullString("R"),
            Chamber:        sqlNullString("House"),
            Symbol:         "XOM",
            Action:         "Buy",
            Amount:         sqlNullString("$50k - $100k"),
            ExecutedAt:     time.Now().AddDate(0, 0, -12),
            DisclosureDate: time.Now().AddDate(0, 0, -6),
            Sentiment:      0.22,
            SourceUrl:      sqlNullString("https://clerk.house.gov"),
        },
    }

    for _, event := range events {
        if err := queries.InsertCongressTrade(ctx, event); err != nil {
            return err
        }
    }

    log.Info("seeded congress trades")
    return nil
}

func seedRecommendations(ctx context.Context, queries *database.Queries, log *slog.Logger, uuidFn func() string) error {
    rows, err := queries.ListRecommendations(ctx, 1)
    if err == nil && len(rows) > 0 {
        return nil
    }

    plays := []database.InsertRecommendationParams{
        {
            ID:         uuidFn(),
            Symbol:     "MSFT",
            Thesis:     "Pair Azure OpenAI uptake with Copilot seat growth for durable 20% cloud ARR.",
            Conviction: "High",
            Score:      8.7,
            Catalyst:   sqlNullString("Ignite + Build launches"),
            CreatedAt:  time.Now().Add(-4 * time.Hour),
        },
        {
            ID:         uuidFn(),
            Symbol:     "NVDA",
            Thesis:     "Next-gen Blackwell supply ramp plus networking upsell sustain EPS beats.",
            Conviction: "Medium",
            Score:      7.9,
            Catalyst:   sqlNullString("Q4 earnings"),
            CreatedAt:  time.Now().Add(-2 * time.Hour),
        },
        {
            ID:         uuidFn(),
            Symbol:     "LMT",
            Thesis:     "Supplemental defense package pushes book-to-bill above 1.3x.",
            Conviction: "Watchlist",
            Score:      6.1,
            Catalyst:   sqlNullString("FY25 NDAA"),
            CreatedAt:  time.Now().Add(-1 * time.Hour),
        },
    }

    for _, play := range plays {
        if err := queries.InsertRecommendation(ctx, play); err != nil {
            return err
        }
    }

    log.Info("seeded recommendations")
    return nil
}

func sqlNullString(value string) sql.NullString {
    if value == "" {
        return sql.NullString{}
    }
    return sql.NullString{String: value, Valid: true}
}
