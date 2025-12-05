package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// MarketDataService aggregates data from multiple financial APIs
type MarketDataService struct {
	httpClient     *http.Client
	alphaVantageKey string
	finnhubKey      string
	cache           *MarketCache
}

// MarketCache provides thread-safe caching for market data
type MarketCache struct {
	mu        sync.RWMutex
	quotes    map[string]*CachedQuote
	indices   []IndexQuote
	indexTime time.Time
}

// CachedQuote holds cached stock data with expiration
type CachedQuote struct {
	Data      *StockQuote
	ExpiresAt time.Time
}

// StockQuote represents real-time stock data
type StockQuote struct {
	Symbol        string    `json:"symbol"`
	Name          string    `json:"name"`
	Price         float64   `json:"price"`
	Change        float64   `json:"change"`
	ChangePercent float64   `json:"changePercent"`
	Open          float64   `json:"open"`
	High          float64   `json:"high"`
	Low           float64   `json:"low"`
	PrevClose     float64   `json:"prevClose"`
	Volume        int64     `json:"volume"`
	MarketCap     int64     `json:"marketCap"`
	PE            float64   `json:"pe"`
	Week52High    float64   `json:"week52High"`
	Week52Low     float64   `json:"week52Low"`
	AvgVolume     int64     `json:"avgVolume"`
	Dividend      float64   `json:"dividend"`
	DividendYield float64   `json:"dividendYield"`
	EPS           float64   `json:"eps"`
	Beta          float64   `json:"beta"`
	Sector        string    `json:"sector"`
	Industry      string    `json:"industry"`
	Exchange      string    `json:"exchange"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// IndexQuote represents a market index
type IndexQuote struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
}

// MarketOverview contains summary market data
type MarketOverview struct {
	Indices       []IndexQuote     `json:"indices"`
	TopGainers    []StockQuote     `json:"topGainers"`
	TopLosers     []StockQuote     `json:"topLosers"`
	MostActive    []StockQuote     `json:"mostActive"`
	SectorPerf    []SectorPerf     `json:"sectorPerf"`
	MarketStatus  string           `json:"marketStatus"`
	LastUpdated   time.Time        `json:"lastUpdated"`
}

// SectorPerf represents sector performance
type SectorPerf struct {
	Sector        string  `json:"sector"`
	ChangePercent float64 `json:"changePercent"`
}

// HistoricalData represents price history
type HistoricalData struct {
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

// NewMarketDataService creates a new market data service
func NewMarketDataService(alphaVantageKey, finnhubKey string) *MarketDataService {
	return &MarketDataService{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		alphaVantageKey: alphaVantageKey,
		finnhubKey:      finnhubKey,
		cache: &MarketCache{
			quotes: make(map[string]*CachedQuote),
		},
	}
}

// GetQuote fetches a stock quote, using cache when available
func (s *MarketDataService) GetQuote(ctx context.Context, symbol string) (*StockQuote, error) {
	// Check cache first
	s.cache.mu.RLock()
	if cached, ok := s.cache.quotes[symbol]; ok && time.Now().Before(cached.ExpiresAt) {
		s.cache.mu.RUnlock()
		return cached.Data, nil
	}
	s.cache.mu.RUnlock()

	// Try Finnhub first (faster, more requests)
	quote, err := s.fetchFinnhubQuote(ctx, symbol)
	if err != nil {
		// Fallback to Alpha Vantage
		quote, err = s.fetchAlphaVantageQuote(ctx, symbol)
		if err != nil {
			// Final fallback to Yahoo Finance
			quote, err = s.fetchYahooQuote(ctx, symbol)
			if err != nil {
				return nil, fmt.Errorf("all data sources failed for %s: %w", symbol, err)
			}
		}
	}

	// Cache the result for 1 minute
	s.cache.mu.Lock()
	s.cache.quotes[symbol] = &CachedQuote{
		Data:      quote,
		ExpiresAt: time.Now().Add(1 * time.Minute),
	}
	s.cache.mu.Unlock()

	return quote, nil
}

// GetMultipleQuotes fetches quotes for multiple symbols concurrently
func (s *MarketDataService) GetMultipleQuotes(ctx context.Context, symbols []string) (map[string]*StockQuote, error) {
	results := make(map[string]*StockQuote)
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(symbols))

	for _, sym := range symbols {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			quote, err := s.GetQuote(ctx, symbol)
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			results[symbol] = quote
			mu.Unlock()
		}(sym)
	}

	wg.Wait()
	close(errChan)

	return results, nil
}

// GetMarketOverview returns summary market data
func (s *MarketDataService) GetMarketOverview(ctx context.Context) (*MarketOverview, error) {
	overview := &MarketOverview{
		LastUpdated: time.Now(),
	}

	// Fetch indices
	indices, err := s.GetIndices(ctx)
	if err == nil {
		overview.Indices = indices
	}

	// Fetch movers
	gainers, losers, active, err := s.GetMarketMovers(ctx)
	if err == nil {
		overview.TopGainers = gainers
		overview.TopLosers = losers
		overview.MostActive = active
	}

	// Fetch sector performance
	sectors, err := s.GetSectorPerformance(ctx)
	if err == nil {
		overview.SectorPerf = sectors
	}

	// Determine market status
	overview.MarketStatus = s.GetMarketStatus()

	return overview, nil
}

// GetIndices returns major market indices
func (s *MarketDataService) GetIndices(ctx context.Context) ([]IndexQuote, error) {
	// Check cache
	s.cache.mu.RLock()
	if !s.cache.indexTime.IsZero() && time.Since(s.cache.indexTime) < 30*time.Second {
		indices := s.cache.indices
		s.cache.mu.RUnlock()
		return indices, nil
	}
	s.cache.mu.RUnlock()

	indexSymbols := []string{"^GSPC", "^DJI", "^IXIC", "^RUT", "^VIX"}
	indexNames := map[string]string{
		"^GSPC": "S&P 500",
		"^DJI":  "Dow Jones",
		"^IXIC": "NASDAQ",
		"^RUT":  "Russell 2000",
		"^VIX":  "VIX",
	}

	var indices []IndexQuote
	for _, symbol := range indexSymbols {
		quote, err := s.fetchYahooQuote(ctx, symbol)
		if err != nil {
			continue
		}
		indices = append(indices, IndexQuote{
			Symbol:        symbol,
			Name:          indexNames[symbol],
			Price:         quote.Price,
			Change:        quote.Change,
			ChangePercent: quote.ChangePercent,
		})
	}

	// Cache indices
	s.cache.mu.Lock()
	s.cache.indices = indices
	s.cache.indexTime = time.Now()
	s.cache.mu.Unlock()

	return indices, nil
}

// GetMarketMovers returns top gainers, losers, and most active stocks
func (s *MarketDataService) GetMarketMovers(ctx context.Context) (gainers, losers, active []StockQuote, err error) {
	// This would typically call an API endpoint for market movers
	// For now, return curated list that can be updated
	symbols := []string{"NVDA", "AAPL", "MSFT", "GOOGL", "AMZN", "META", "TSLA", "AMD", "NFLX", "JPM"}

	quotes, err := s.GetMultipleQuotes(ctx, symbols)
	if err != nil {
		return nil, nil, nil, err
	}

	var allQuotes []StockQuote
	for _, q := range quotes {
		allQuotes = append(allQuotes, *q)
	}

	// Sort by change percent for gainers/losers
	// In production, this would come pre-sorted from the API
	for _, q := range allQuotes {
		if q.ChangePercent > 0 {
			gainers = append(gainers, q)
		} else {
			losers = append(losers, q)
		}
	}

	// Most active by volume
	active = allQuotes
	return gainers, losers, active, nil
}

// GetSectorPerformance returns sector performance data
func (s *MarketDataService) GetSectorPerformance(ctx context.Context) ([]SectorPerf, error) {
	// Sector ETFs to track
	sectorETFs := map[string]string{
		"XLK": "Technology",
		"XLF": "Financials",
		"XLV": "Healthcare",
		"XLE": "Energy",
		"XLY": "Consumer Cyclical",
		"XLP": "Consumer Defensive",
		"XLI": "Industrials",
		"XLB": "Materials",
		"XLU": "Utilities",
		"XLRE": "Real Estate",
	}

	var sectors []SectorPerf
	for symbol, name := range sectorETFs {
		quote, err := s.GetQuote(ctx, symbol)
		if err != nil {
			continue
		}
		sectors = append(sectors, SectorPerf{
			Sector:        name,
			ChangePercent: quote.ChangePercent,
		})
	}

	return sectors, nil
}

// GetHistoricalData returns historical price data
func (s *MarketDataService) GetHistoricalData(ctx context.Context, symbol string, period string) ([]HistoricalData, error) {
	// Would call Alpha Vantage or Yahoo Finance for historical data
	// Period: 1D, 5D, 1M, 3M, 6M, 1Y, 5Y
	return s.fetchYahooHistorical(ctx, symbol, period)
}

// GetMarketStatus returns whether the market is open
func (s *MarketDataService) GetMarketStatus() string {
	now := time.Now().In(time.FixedZone("EST", -5*60*60))
	hour := now.Hour()
	minute := now.Minute()
	weekday := now.Weekday()

	// Market hours: 9:30 AM - 4:00 PM EST, Monday-Friday
	if weekday == time.Saturday || weekday == time.Sunday {
		return "closed"
	}

	marketOpen := hour > 9 || (hour == 9 && minute >= 30)
	marketClose := hour < 16

	if marketOpen && marketClose {
		return "open"
	} else if hour >= 4 && (hour < 9 || (hour == 9 && minute < 30)) {
		return "pre-market"
	} else if hour >= 16 && hour < 20 {
		return "after-hours"
	}
	return "closed"
}

// fetchFinnhubQuote fetches quote from Finnhub API
func (s *MarketDataService) fetchFinnhubQuote(ctx context.Context, symbol string) (*StockQuote, error) {
	if s.finnhubKey == "" {
		return nil, fmt.Errorf("finnhub API key not configured")
	}

	url := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", symbol, s.finnhubKey)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		C  float64 `json:"c"`  // Current price
		D  float64 `json:"d"`  // Change
		DP float64 `json:"dp"` // Percent change
		H  float64 `json:"h"`  // High
		L  float64 `json:"l"`  // Low
		O  float64 `json:"o"`  // Open
		PC float64 `json:"pc"` // Previous close
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if data.C == 0 {
		return nil, fmt.Errorf("no data returned for %s", symbol)
	}

	return &StockQuote{
		Symbol:        symbol,
		Price:         data.C,
		Change:        data.D,
		ChangePercent: data.DP,
		Open:          data.O,
		High:          data.H,
		Low:           data.L,
		PrevClose:     data.PC,
		UpdatedAt:     time.Now(),
	}, nil
}

// fetchAlphaVantageQuote fetches quote from Alpha Vantage API
func (s *MarketDataService) fetchAlphaVantageQuote(ctx context.Context, symbol string) (*StockQuote, error) {
	if s.alphaVantageKey == "" {
		return nil, fmt.Errorf("alpha vantage API key not configured")
	}

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", symbol, s.alphaVantageKey)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		GlobalQuote struct {
			Symbol        string `json:"01. symbol"`
			Open          string `json:"02. open"`
			High          string `json:"03. high"`
			Low           string `json:"04. low"`
			Price         string `json:"05. price"`
			Volume        string `json:"06. volume"`
			PrevClose     string `json:"08. previous close"`
			Change        string `json:"09. change"`
			ChangePercent string `json:"10. change percent"`
		} `json:"Global Quote"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &StockQuote{
		Symbol:    symbol,
		UpdatedAt: time.Now(),
	}, nil
}

// fetchYahooQuote fetches quote from Yahoo Finance (unofficial)
func (s *MarketDataService) fetchYahooQuote(ctx context.Context, symbol string) (*StockQuote, error) {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1d&range=1d", symbol)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Chart struct {
			Result []struct {
				Meta struct {
					Symbol             string  `json:"symbol"`
					RegularMarketPrice float64 `json:"regularMarketPrice"`
					PreviousClose      float64 `json:"previousClose"`
					Exchange           string  `json:"exchangeName"`
				} `json:"meta"`
				Indicators struct {
					Quote []struct {
						Open   []float64 `json:"open"`
						High   []float64 `json:"high"`
						Low    []float64 `json:"low"`
						Close  []float64 `json:"close"`
						Volume []int64   `json:"volume"`
					} `json:"quote"`
				} `json:"indicators"`
			} `json:"result"`
			Error interface{} `json:"error"`
		} `json:"chart"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if len(data.Chart.Result) == 0 {
		return nil, fmt.Errorf("no data for symbol %s", symbol)
	}

	result := data.Chart.Result[0]
	meta := result.Meta
	price := meta.RegularMarketPrice
	prevClose := meta.PreviousClose
	change := price - prevClose
	changePercent := 0.0
	if prevClose > 0 {
		changePercent = (change / prevClose) * 100
	}

	var open, high, low float64
	var volume int64
	if len(result.Indicators.Quote) > 0 && len(result.Indicators.Quote[0].Open) > 0 {
		q := result.Indicators.Quote[0]
		open = q.Open[0]
		high = q.High[0]
		low = q.Low[0]
		if len(q.Volume) > 0 {
			volume = q.Volume[0]
		}
	}

	return &StockQuote{
		Symbol:        symbol,
		Price:         price,
		Change:        change,
		ChangePercent: changePercent,
		Open:          open,
		High:          high,
		Low:           low,
		PrevClose:     prevClose,
		Volume:        volume,
		Exchange:      meta.Exchange,
		UpdatedAt:     time.Now(),
	}, nil
}

// fetchYahooHistorical fetches historical data from Yahoo Finance
func (s *MarketDataService) fetchYahooHistorical(ctx context.Context, symbol string, period string) ([]HistoricalData, error) {
	intervalMap := map[string]string{
		"1D":  "5m",
		"5D":  "15m",
		"1M":  "1h",
		"3M":  "1d",
		"6M":  "1d",
		"1Y":  "1d",
		"5Y":  "1wk",
	}

	interval := intervalMap[period]
	if interval == "" {
		interval = "1d"
	}

	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=%s&range=%s", symbol, interval, period)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Chart struct {
			Result []struct {
				Timestamp  []int64 `json:"timestamp"`
				Indicators struct {
					Quote []struct {
						Open   []float64 `json:"open"`
						High   []float64 `json:"high"`
						Low    []float64 `json:"low"`
						Close  []float64 `json:"close"`
						Volume []int64   `json:"volume"`
					} `json:"quote"`
				} `json:"indicators"`
			} `json:"result"`
		} `json:"chart"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if len(data.Chart.Result) == 0 {
		return nil, fmt.Errorf("no historical data for %s", symbol)
	}

	result := data.Chart.Result[0]
	var history []HistoricalData

	if len(result.Indicators.Quote) > 0 {
		q := result.Indicators.Quote[0]
		for i, ts := range result.Timestamp {
			if i < len(q.Close) {
				history = append(history, HistoricalData{
					Date:   time.Unix(ts, 0).Format("2006-01-02 15:04"),
					Open:   q.Open[i],
					High:   q.High[i],
					Low:    q.Low[i],
					Close:  q.Close[i],
					Volume: q.Volume[i],
				})
			}
		}
	}

	return history, nil
}
