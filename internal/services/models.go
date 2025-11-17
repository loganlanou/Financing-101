package services

import "time"

// NewsHeadline captures the NLP feature-set for a scraped article.
type NewsHeadline struct {
    ID            string
    Title         string
    Source        string
    Summary       string
    Sentiment     float64
    Trend         string
    Tickers       []string
    URL           string
    PublishedAt   time.Time
}

// StockSnapshot compares a holding vs S&P500 across timeframes.
type StockSnapshot struct {
    ID           string
    Symbol       string
    Name         string
    Sector       string
    Industry     string
    Change30     float64
    Change90     float64
    Change365    float64
    VsSP500_30   float64
    VsSP500_90   float64
    VsSP500_365  float64
    Conviction   string
    Thesis       string
    UpdatedAt    time.Time
}

// Trade mirrors congressional disclosure data.
type Trade struct {
    ID             string
    Member         string
    Party          string
    Chamber        string
    Symbol         string
    Action         string
    Amount         string
    ExecutedAt     time.Time
    DisclosureDate time.Time
    Sentiment      float64
    SourceURL      string
}

// Recommendation is a curated insight published on the dashboard.
type Recommendation struct {
    ID         string
    Symbol     string
    Thesis     string
    Conviction string
    Score      float64
    Catalyst   string
    CreatedAt  time.Time
}

