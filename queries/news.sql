-- name: ListLatestNews :many
SELECT id, title, source, summary, sentiment_score, trend, tickers, url, published_at
FROM news_articles
ORDER BY published_at DESC
LIMIT sqlc.arg('limit');

-- name: InsertNewsArticle :exec
INSERT INTO news_articles (id, title, source, summary, sentiment_score, trend, tickers, url, published_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    title=excluded.title,
    source=excluded.source,
    summary=excluded.summary,
    sentiment_score=excluded.sentiment_score,
    trend=excluded.trend,
    tickers=excluded.tickers,
    url=excluded.url,
    published_at=excluded.published_at;
