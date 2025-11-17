-- name: ListRecentTrades :many
SELECT id, member, party, chamber, symbol, action, amount, executed_at, disclosure_date, sentiment, source_url
FROM congress_trades
ORDER BY executed_at DESC
LIMIT sqlc.arg('limit');

-- name: InsertCongressTrade :exec
INSERT INTO congress_trades (
    id, member, party, chamber, symbol, action, amount, executed_at, disclosure_date, sentiment, source_url
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    member=excluded.member,
    party=excluded.party,
    chamber=excluded.chamber,
    symbol=excluded.symbol,
    action=excluded.action,
    amount=excluded.amount,
    executed_at=excluded.executed_at,
    disclosure_date=excluded.disclosure_date,
    sentiment=excluded.sentiment,
    source_url=excluded.source_url;
