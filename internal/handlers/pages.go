package handlers

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/loganlanou/Financing-101/internal/services"
	"github.com/loganlanou/Financing-101/web/components/pages"
	"golang.org/x/sync/errgroup"
	"log/slog"
)

// PagesHandler serves all page routes with the modern Layout
type PagesHandler struct {
	log          *slog.Logger
	newsService  *services.NewsService
	stockService *services.StockService
	tradeService *services.TradeService
	recService   *services.RecommendationService
}

func NewPagesHandler(
	log *slog.Logger,
	newsService *services.NewsService,
	stockService *services.StockService,
	tradeService *services.TradeService,
	recService *services.RecommendationService,
) *PagesHandler {
	return &PagesHandler{
		log:          log,
		newsService:  newsService,
		stockService: stockService,
		tradeService: tradeService,
		recService:   recService,
	}
}

func (h *PagesHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/", h.dashboard)
	e.GET("/markets", h.markets)
	e.GET("/stocks", h.stocks)
}

func (h *PagesHandler) dashboard(c echo.Context) error {
	reqCtx := c.Request().Context()

	var (
		news   []services.NewsHeadline
		trades []services.Trade
		recs   []services.Recommendation
	)

	g, ctx := errgroup.WithContext(reqCtx)

	g.Go(func() error {
		data, err := h.newsService.Latest(ctx, 8)
		if err != nil {
			return err
		}
		news = data
		return nil
	})

	g.Go(func() error {
		data, err := h.tradeService.Recent(ctx, 6)
		if err != nil {
			return err
		}
		trades = data
		return nil
	})

	g.Go(func() error {
		data, err := h.recService.TopPicks(ctx, 4)
		if err != nil {
			return err
		}
		recs = data
		return nil
	})

	if err := g.Wait(); err != nil {
		h.log.Error("dashboard aggregation failed", slog.Any("err", err))
	}

	indices := getMockIndices()
	movers := getMockGainersLosers()
	marketStatus := getMarketStatus()

	data := pages.DashboardData{
		Indices:         indices,
		TopGainers:      movers.gainers,
		TopLosers:       movers.losers,
		RecentNews:      news,
		CongressTrades:  trades,
		Recommendations: recs,
		MarketStatus:    marketStatus,
		LastUpdated:     time.Now(),
	}

	page := pages.DashboardPage(data)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return page.Render(reqCtx, c.Response())
}

func (h *PagesHandler) markets(c echo.Context) error {
	reqCtx := c.Request().Context()

	indices := getMockIndices()
	sectors := getMockSectors()
	movers := getMockGainersLosers()

	data := pages.MarketsData{
		Indices:      indices,
		Sectors:      sectors,
		TopGainers:   movers.gainers,
		TopLosers:    movers.losers,
		MostActive:   movers.active,
		MarketStatus: getMarketStatus(),
	}

	page := pages.MarketsPage(data)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return page.Render(reqCtx, c.Response())
}

func (h *PagesHandler) stocks(c echo.Context) error {
	reqCtx := c.Request().Context()

	stocks := getMockStocks()
	var featured *services.StockQuote
	if len(stocks) > 0 {
		featured = &stocks[0]
	}

	data := pages.StocksData{
		Stocks:        stocks,
		FeaturedStock: featured,
	}

	page := pages.StocksPage(data)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return page.Render(reqCtx, c.Response())
}

// Helper functions

func getMarketStatus() string {
	now := time.Now()
	hour := now.Hour()
	weekday := now.Weekday()

	if weekday == time.Saturday || weekday == time.Sunday {
		return "closed"
	}
	if hour >= 9 && hour < 16 {
		return "open"
	}
	if hour >= 4 && hour < 9 {
		return "pre-market"
	}
	if hour >= 16 && hour < 20 {
		return "after-hours"
	}
	return "closed"
}

func getMockIndices() []services.IndexQuote {
	return []services.IndexQuote{
		{Symbol: "^GSPC", Name: "S&P 500", Price: 5998.74, Change: 63.77, ChangePercent: 1.07},
		{Symbol: "^DJI", Name: "Dow Jones", Price: 43828.06, Change: 619.05, ChangePercent: 1.43},
		{Symbol: "^IXIC", Name: "NASDAQ", Price: 19926.72, Change: 180.09, ChangePercent: 0.91},
		{Symbol: "^RUT", Name: "Russell 2000", Price: 2346.90, Change: 51.27, ChangePercent: 2.23},
	}
}

type mockMovers struct {
	gainers []services.StockQuote
	losers  []services.StockQuote
	active  []services.StockQuote
}

func getMockGainersLosers() mockMovers {
	gainers := []services.StockQuote{
		{Symbol: "NVDA", Name: "NVIDIA Corporation", Price: 134.25, Change: 8.42, ChangePercent: 6.69, Volume: 312500000, MarketCap: 3300000000000, PE: 65.2, Week52High: 140.76, Week52Low: 45.01},
		{Symbol: "AMD", Name: "Advanced Micro Devices", Price: 137.89, Change: 5.23, ChangePercent: 3.94, Volume: 48200000, MarketCap: 223000000000, PE: 98.5, Week52High: 164.46, Week52Low: 93.12},
		{Symbol: "TSLA", Name: "Tesla, Inc.", Price: 424.77, Change: 14.89, ChangePercent: 3.63, Volume: 89500000, MarketCap: 1350000000000, PE: 78.4, Week52High: 438.22, Week52Low: 138.80},
		{Symbol: "GOOGL", Name: "Alphabet Inc.", Price: 192.96, Change: 5.67, ChangePercent: 3.03, Volume: 28400000, MarketCap: 2370000000000, PE: 24.8, Week52High: 193.31, Week52Low: 130.67},
		{Symbol: "META", Name: "Meta Platforms", Price: 617.12, Change: 15.23, ChangePercent: 2.53, Volume: 14200000, MarketCap: 1560000000000, PE: 29.3, Week52High: 618.95, Week52Low: 326.89},
	}

	losers := []services.StockQuote{
		{Symbol: "INTC", Name: "Intel Corporation", Price: 20.11, Change: -0.89, ChangePercent: -4.24, Volume: 67800000, MarketCap: 86000000000, PE: 0, Week52High: 51.28, Week52Low: 18.51},
		{Symbol: "BA", Name: "Boeing Company", Price: 177.56, Change: -4.23, ChangePercent: -2.33, Volume: 8900000, MarketCap: 108000000000, PE: 0, Week52High: 267.54, Week52Low: 137.03},
		{Symbol: "NKE", Name: "Nike, Inc.", Price: 75.89, Change: -1.67, ChangePercent: -2.15, Volume: 12300000, MarketCap: 113000000000, PE: 21.4, Week52High: 107.43, Week52Low: 70.75},
		{Symbol: "DIS", Name: "Walt Disney Co.", Price: 112.34, Change: -2.12, ChangePercent: -1.85, Volume: 9800000, MarketCap: 205000000000, PE: 72.1, Week52High: 123.74, Week52Low: 83.91},
		{Symbol: "PFE", Name: "Pfizer Inc.", Price: 25.67, Change: -0.43, ChangePercent: -1.65, Volume: 32100000, MarketCap: 145000000000, PE: 19.8, Week52High: 31.54, Week52Low: 24.48},
	}

	active := []services.StockQuote{
		{Symbol: "NVDA", Name: "NVIDIA Corporation", Price: 134.25, Change: 8.42, ChangePercent: 6.69, Volume: 312500000, MarketCap: 3300000000000, PE: 65.2, Week52High: 140.76, Week52Low: 45.01},
		{Symbol: "TSLA", Name: "Tesla, Inc.", Price: 424.77, Change: 14.89, ChangePercent: 3.63, Volume: 89500000, MarketCap: 1350000000000, PE: 78.4, Week52High: 438.22, Week52Low: 138.80},
		{Symbol: "INTC", Name: "Intel Corporation", Price: 20.11, Change: -0.89, ChangePercent: -4.24, Volume: 67800000, MarketCap: 86000000000, PE: 0, Week52High: 51.28, Week52Low: 18.51},
		{Symbol: "AMD", Name: "Advanced Micro Devices", Price: 137.89, Change: 5.23, ChangePercent: 3.94, Volume: 48200000, MarketCap: 223000000000, PE: 98.5, Week52High: 164.46, Week52Low: 93.12},
		{Symbol: "AAPL", Name: "Apple Inc.", Price: 248.13, Change: 2.87, ChangePercent: 1.17, Volume: 45600000, MarketCap: 3780000000000, PE: 32.1, Week52High: 250.00, Week52Low: 164.08},
	}

	return mockMovers{gainers: gainers, losers: losers, active: active}
}

func getMockSectors() []services.SectorPerf {
	return []services.SectorPerf{
		{Sector: "Technology", ChangePercent: 2.34},
		{Sector: "Healthcare", ChangePercent: 1.12},
		{Sector: "Financials", ChangePercent: 0.89},
		{Sector: "Energy", ChangePercent: -0.45},
		{Sector: "Consumer Cyclical", ChangePercent: 1.67},
		{Sector: "Consumer Defensive", ChangePercent: 0.23},
		{Sector: "Industrials", ChangePercent: 0.78},
		{Sector: "Materials", ChangePercent: -0.12},
		{Sector: "Utilities", ChangePercent: -0.34},
		{Sector: "Real Estate", ChangePercent: 0.56},
	}
}

func getMockStocks() []services.StockQuote {
	return []services.StockQuote{
		{Symbol: "AAPL", Name: "Apple Inc.", Price: 248.13, Change: 2.87, ChangePercent: 1.17, Open: 245.50, High: 249.25, Low: 244.80, PrevClose: 245.26, Volume: 45600000, MarketCap: 3780000000000, PE: 32.1, Week52High: 250.00, Week52Low: 164.08},
		{Symbol: "MSFT", Name: "Microsoft Corporation", Price: 446.95, Change: 3.12, ChangePercent: 0.70, Open: 443.50, High: 448.75, Low: 442.20, PrevClose: 443.83, Volume: 21400000, MarketCap: 3320000000000, PE: 36.8, Week52High: 468.35, Week52Low: 362.90},
		{Symbol: "NVDA", Name: "NVIDIA Corporation", Price: 134.25, Change: 8.42, ChangePercent: 6.69, Open: 126.50, High: 135.50, Low: 125.80, PrevClose: 125.83, Volume: 312500000, MarketCap: 3300000000000, PE: 65.2, Week52High: 140.76, Week52Low: 45.01},
		{Symbol: "GOOGL", Name: "Alphabet Inc.", Price: 192.96, Change: 5.67, ChangePercent: 3.03, Open: 188.00, High: 193.50, Low: 187.25, PrevClose: 187.29, Volume: 28400000, MarketCap: 2370000000000, PE: 24.8, Week52High: 193.31, Week52Low: 130.67},
		{Symbol: "AMZN", Name: "Amazon.com, Inc.", Price: 227.03, Change: 4.23, ChangePercent: 1.90, Open: 223.50, High: 228.00, Low: 222.80, PrevClose: 222.80, Volume: 39800000, MarketCap: 2380000000000, PE: 46.2, Week52High: 233.00, Week52Low: 151.61},
		{Symbol: "META", Name: "Meta Platforms", Price: 617.12, Change: 15.23, ChangePercent: 2.53, Open: 602.00, High: 618.95, Low: 600.50, PrevClose: 601.89, Volume: 14200000, MarketCap: 1560000000000, PE: 29.3, Week52High: 618.95, Week52Low: 326.89},
		{Symbol: "TSLA", Name: "Tesla, Inc.", Price: 424.77, Change: 14.89, ChangePercent: 3.63, Open: 410.00, High: 426.50, Low: 408.25, PrevClose: 409.88, Volume: 89500000, MarketCap: 1350000000000, PE: 78.4, Week52High: 438.22, Week52Low: 138.80},
		{Symbol: "BRK.B", Name: "Berkshire Hathaway", Price: 458.92, Change: -1.23, ChangePercent: -0.27, Open: 460.00, High: 461.50, Low: 457.80, PrevClose: 460.15, Volume: 3200000, MarketCap: 989000000000, PE: 9.8, Week52High: 491.66, Week52Low: 378.94},
		{Symbol: "JPM", Name: "JPMorgan Chase", Price: 243.67, Change: 3.45, ChangePercent: 1.44, Open: 240.50, High: 244.50, Low: 239.80, PrevClose: 240.22, Volume: 8900000, MarketCap: 698000000000, PE: 12.5, Week52High: 254.31, Week52Low: 173.21},
		{Symbol: "V", Name: "Visa Inc.", Price: 317.89, Change: 2.34, ChangePercent: 0.74, Open: 315.50, High: 318.75, Low: 314.80, PrevClose: 315.55, Volume: 6200000, MarketCap: 628000000000, PE: 30.2, Week52High: 321.62, Week52Low: 252.70},
	}
}
