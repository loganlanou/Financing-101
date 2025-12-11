package handlers

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/loganlanou/Financing-101/internal/services"
    "github.com/loganlanou/Financing-101/web/components"
    "golang.org/x/sync/errgroup"
    "log/slog"
)

// DashboardHandler wires the core Investing101 experiences.
type DashboardHandler struct {
    log             *slog.Logger
    newsService     *services.NewsService
    stockService    *services.StockService
    tradeService    *services.TradeService
    recService      *services.RecommendationService
}

func NewDashboardHandler(
    log *slog.Logger,
    newsService *services.NewsService,
    stockService *services.StockService,
    tradeService *services.TradeService,
    recService *services.RecommendationService,
) *DashboardHandler {
    return &DashboardHandler{
        log:          log,
        newsService:  newsService,
        stockService: stockService,
        tradeService: tradeService,
        recService:   recService,
    }
}

func (h *DashboardHandler) RegisterRoutes(e *echo.Echo) {
    e.GET("/", h.dashboard)
}

func (h *DashboardHandler) dashboard(c echo.Context) error {
    reqCtx := c.Request().Context()
    var (
        news []services.NewsHeadline
        stocks []services.StockSnapshot
        trades []services.Trade
        recs []services.Recommendation
    )

    g, ctx := errgroup.WithContext(reqCtx)

    g.Go(func() error {
        data, err := h.newsService.Latest(ctx, 6)
        if err != nil {
            return err
        }
        news = data
        return nil
    })

    g.Go(func() error {
        data, err := h.stockService.Leaders(ctx, 6)
        if err != nil {
            return err
        }
        stocks = data
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
        return c.String(http.StatusInternalServerError, "unable to load dashboard")
    }

    metrics := components.DashboardMetrics{
        NewsCount:           len(news),
        StocksOutperforming: countOutperforming(stocks),
        CongressionalTrades: len(trades),
        AvgConviction:       avgScore(recs),
    }

    page := components.DashboardPage(components.DashboardPageData{
        News:            news,
        Stocks:          stocks,
        Trades:          trades,
        Recommendations: recs,
        Metrics:         metrics,
    })

    c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
    if err := page.Render(reqCtx, c.Response()); err != nil {
        h.log.Error("template render failed", slog.Any("err", err))
        return err
    }

    return nil
}

func countOutperforming(stocks []services.StockSnapshot) int {
    count := 0
    for _, stock := range stocks {
        if stock.VsSP500_90 > 0 {
            count++
        }
    }

    return count
}

func avgScore(recs []services.Recommendation) float64 {
    if len(recs) == 0 {
        return 0
    }

    total := 0.0
    for _, rec := range recs {
        total += rec.Score
    }

    return total / float64(len(recs))
}

