package services

import (
    "context"

    "github.com/loganlanou/Financing-101/internal/database"
    "log/slog"
)

// TradeService exposes congressional trading data.
type TradeService struct {
    log     *slog.Logger
    queries *database.Queries
}

func NewTradeService(log *slog.Logger, queries *database.Queries) *TradeService {
    return &TradeService{log: log, queries: queries}
}

func (s *TradeService) Recent(ctx context.Context, limit int32) ([]Trade, error) {
    rows, err := s.queries.ListRecentTrades(ctx, int64(limit))
    if err != nil {
        return nil, err
    }

    trades := make([]Trade, 0, len(rows))
    for _, row := range rows {
        trades = append(trades, Trade{
            ID:             row.ID,
            Member:         row.Member,
            Party:          row.Party.String,
            Chamber:        row.Chamber.String,
            Symbol:         row.Symbol,
            Action:         row.Action,
            Amount:         row.Amount.String,
            ExecutedAt:     row.ExecutedAt,
            DisclosureDate: row.DisclosureDate,
            Sentiment:      row.Sentiment,
            SourceURL:      row.SourceUrl.String,
        })
    }

    return trades, nil
}
