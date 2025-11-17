-- +goose Up
CREATE TABLE IF NOT EXISTS news_articles (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    source TEXT NOT NULL,
    summary TEXT NOT NULL,
    sentiment_score REAL NOT NULL,
    trend TEXT NOT NULL,
    tickers TEXT NOT NULL,
    url TEXT NOT NULL,
    published_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS stock_snapshots (
    id TEXT PRIMARY KEY,
    symbol TEXT NOT NULL,
    name TEXT NOT NULL,
    sector TEXT,
    industry TEXT,
    change_30 REAL NOT NULL,
    change_90 REAL NOT NULL,
    change_365 REAL NOT NULL,
    vs_sp500_30 REAL NOT NULL,
    vs_sp500_90 REAL NOT NULL,
    vs_sp500_365 REAL NOT NULL,
    conviction TEXT NOT NULL,
    thesis TEXT NOT NULL,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS congress_trades (
    id TEXT PRIMARY KEY,
    member TEXT NOT NULL,
    party TEXT,
    chamber TEXT,
    symbol TEXT NOT NULL,
    action TEXT NOT NULL,
    amount TEXT,
    executed_at DATETIME NOT NULL,
    disclosure_date DATETIME NOT NULL,
    sentiment REAL NOT NULL DEFAULT 0,
    source_url TEXT
);

CREATE TABLE IF NOT EXISTS recommendations (
    id TEXT PRIMARY KEY,
    symbol TEXT NOT NULL,
    thesis TEXT NOT NULL,
    conviction TEXT NOT NULL,
    score REAL NOT NULL,
    catalyst TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS recommendations;
DROP TABLE IF EXISTS congress_trades;
DROP TABLE IF EXISTS stock_snapshots;
DROP TABLE IF EXISTS news_articles;
