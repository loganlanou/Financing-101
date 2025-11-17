-- name: ListRecommendations :many
SELECT id, symbol, thesis, conviction, score, catalyst, created_at
FROM recommendations
ORDER BY score DESC
LIMIT sqlc.arg('limit');

-- name: InsertRecommendation :exec
INSERT INTO recommendations (id, symbol, thesis, conviction, score, catalyst, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    symbol=excluded.symbol,
    thesis=excluded.thesis,
    conviction=excluded.conviction,
    score=excluded.score,
    catalyst=excluded.catalyst,
    created_at=excluded.created_at;
