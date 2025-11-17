package services

import (
    "context"

    "github.com/loganlanou/Financing-101/internal/database"
    "log/slog"
)

// StockService surfaces performance vs benchmarks.
type StockService struct {
    log     *slog.Logger
    queries *database.Queries
}

func NewStockService(log *slog.Logger, queries *database.Queries) *StockService {
    return &StockService{log: log, queries: queries}
}

func (s *StockService) Leaders(ctx context.Context, limit int32) ([]StockSnapshot, error) {
    rows, err := s.queries.ListStockSnapshots(ctx, int64(limit))
    if err != nil {
        return nil, err
    }

    out := make([]StockSnapshot, 0, len(rows))
    for _, row := range rows {
        out = append(out, StockSnapshot{
            ID:          row.ID,
            Symbol:      row.Symbol,
            Name:        row.Name,
            Sector:      row.Sector.String,
            Industry:    row.Industry.String,
            Change30:    row.Change30,
            Change90:    row.Change90,
            Change365:   row.Change365,
            VsSP500_30:  row.VsSp50030,
            VsSP500_90:  row.VsSp50090,
            VsSP500_365: row.VsSp500365,
            Conviction:  row.Conviction,
            Thesis:      row.Thesis,
            UpdatedAt:   row.UpdatedAt,
        })
    }

    return out, nil
}
