-- name: ListStockSnapshots :many
SELECT id, symbol, name, sector, industry, change_30, change_90, change_365,
       vs_sp500_30, vs_sp500_90, vs_sp500_365, conviction, thesis, updated_at
FROM stock_snapshots
ORDER BY vs_sp500_90 DESC
LIMIT sqlc.arg('limit');

-- name: InsertStockSnapshot :exec
INSERT INTO stock_snapshots (
    id, symbol, name, sector, industry, change_30, change_90, change_365,
    vs_sp500_30, vs_sp500_90, vs_sp500_365, conviction, thesis, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    symbol=excluded.symbol,
    name=excluded.name,
    sector=excluded.sector,
    industry=excluded.industry,
    change_30=excluded.change_30,
    change_90=excluded.change_90,
    change_365=excluded.change_365,
    vs_sp500_30=excluded.vs_sp500_30,
    vs_sp500_90=excluded.vs_sp500_90,
    vs_sp500_365=excluded.vs_sp500_365,
    conviction=excluded.conviction,
    thesis=excluded.thesis,
    updated_at=excluded.updated_at;
