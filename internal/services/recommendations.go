package services

import (
    "context"

    "github.com/loganlanou/Financing-101/internal/database"
    "log/slog"
)

// RecommendationService combines news + quant screens for insights.
type RecommendationService struct {
    log     *slog.Logger
    queries *database.Queries
}

func NewRecommendationService(log *slog.Logger, queries *database.Queries) *RecommendationService {
    return &RecommendationService{log: log, queries: queries}
}

func (s *RecommendationService) TopPicks(ctx context.Context, limit int32) ([]Recommendation, error) {
    rows, err := s.queries.ListRecommendations(ctx, int64(limit))
    if err != nil {
        return nil, err
    }

    recs := make([]Recommendation, 0, len(rows))
    for _, row := range rows {
        recs = append(recs, Recommendation{
            ID:         row.ID,
            Symbol:     row.Symbol,
            Thesis:     row.Thesis,
            Conviction: row.Conviction,
            Score:      row.Score,
            Catalyst:   row.Catalyst.String,
            CreatedAt:  row.CreatedAt,
        })
    }

    return recs, nil
}
