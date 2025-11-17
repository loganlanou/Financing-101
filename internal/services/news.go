package services

import (
    "context"
    "strings"

    "github.com/loganlanou/Financing-101/internal/database"
    "log/slog"
)

// NewsService orchestrates NLP-backed article data.
type NewsService struct {
    log     *slog.Logger
    queries *database.Queries
}

func NewNewsService(log *slog.Logger, queries *database.Queries) *NewsService {
    return &NewsService{log: log, queries: queries}
}

func (s *NewsService) Latest(ctx context.Context, limit int32) ([]NewsHeadline, error) {
    rows, err := s.queries.ListLatestNews(ctx, int64(limit))
    if err != nil {
        return nil, err
    }

    headlines := make([]NewsHeadline, 0, len(rows))
    for _, row := range rows {
        headlines = append(headlines, NewsHeadline{
            ID:          row.ID,
            Title:       row.Title,
            Source:      row.Source,
            Summary:     row.Summary,
            Sentiment:   row.SentimentScore,
            Trend:       row.Trend,
            Tickers:     splitTickers(row.Tickers),
            URL:         row.Url,
            PublishedAt: row.PublishedAt,
        })
    }

    return headlines, nil
}

func splitTickers(raw string) []string {
    out := []string{}
    for _, token := range strings.Split(raw, ",") {
        trimmed := strings.TrimSpace(token)
        if trimmed == "" {
            continue
        }
        out = append(out, strings.ToUpper(trimmed))
    }

    return out
}
