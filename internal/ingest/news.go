package ingest

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jonreiter/govader"
	"github.com/loganlanou/Financing-101/internal/database"
	"github.com/mmcdole/gofeed"
	"log/slog"
)

// NewsIngestor pulls RSS feeds and maps them into persisted articles with sentiment.
type NewsIngestor struct {
	log       *slog.Logger
	queries   *database.Queries
	feeds     []string
	parser    *gofeed.Parser
	analyzer  *govader.SentimentIntensityAnalyzer
	tickerLex map[string]struct{}
}

func NewNewsIngestor(log *slog.Logger, queries *database.Queries, feeds []string) *NewsIngestor {
	if len(feeds) == 0 {
		feeds = []string{"https://finance.yahoo.com/news/rssindex"}
	}

	lex := map[string]struct{}{}
	for _, symbol := range []string{"AAPL", "GOOGL", "META", "MSFT", "NVDA", "AMZN", "TSLA", "SPY", "QQQ", "LMT", "XOM", "NFLX", "ORCL", "AMD", "INTC", "AVGO", "JPM", "BAC"} {
		lex[symbol] = struct{}{}
	}

	return &NewsIngestor{
		log:       log,
		queries:   queries,
		feeds:     feeds,
		parser:    gofeed.NewParser(),
		analyzer:  govader.NewSentimentIntensityAnalyzer(),
		tickerLex: lex,
	}
}

// Refresh downloads the feeds and upserts the freshest articles.
func (n *NewsIngestor) Refresh(ctx context.Context, maxArticles int) error {
	if maxArticles <= 0 {
		return errors.New("maxArticles must be positive")
	}

	type article struct {
		title   string
		source  string
		summary string
		link    string
		date    time.Time
	}

	articles := make([]article, 0, maxArticles*len(n.feeds))
	for _, feedURL := range n.feeds {
		feed, err := n.parser.ParseURLWithContext(feedURL, ctx)
		if err != nil {
			n.log.Warn("rss fetch failed", slog.String("feed", feedURL), slog.Any("err", err))
			continue
		}
		for _, item := range feed.Items {
			published := time.Now()
			if item.PublishedParsed != nil {
				published = *item.PublishedParsed
			}
			summary := item.Description
			if summary == "" {
				summary = item.Content
			}
			articles = append(articles, article{
				title:   item.Title,
				source:  feed.Title,
				summary: summary,
				link:    item.Link,
				date:    published,
			})
		}
	}

	if len(articles) == 0 {
		return errors.New("no articles fetched")
	}

	sort.SliceStable(articles, func(i, j int) bool {
		return articles[i].date.After(articles[j].date)
	})

	limit := maxArticles
	if len(articles) < limit {
		limit = len(articles)
	}

	for i := 0; i < limit; i++ {
		art := articles[i]
		score := n.analyzer.PolarityScores(art.title + " " + art.summary).Compound
		trend := sentimentLabel(score)
		tickers := strings.Join(n.extractTickers(art.title+" "+art.summary), ", ")

		arg := database.InsertNewsArticleParams{
			ID:             deterministicID(art.link, art.title, art.date),
			Title:          art.title,
			Source:         art.source,
			Summary:        truncate(art.summary, 280),
			SentimentScore: score,
			Trend:          trend,
			Tickers:        tickers,
			Url:            art.link,
			PublishedAt:    art.date,
		}

		if err := n.queries.InsertNewsArticle(ctx, arg); err != nil {
			n.log.Warn("upsert article failed", slog.String("title", art.title), slog.Any("err", err))
		}
	}

	n.log.Info("news refresh complete", slog.Int("articles", limit))
	return nil
}

func (n *NewsIngestor) extractTickers(text string) []string {
	tokens := strings.FieldsFunc(text, func(r rune) bool {
		switch r {
		case ' ', ',', '.', ';', ':', '(', ')', '"', '\'', '?', '!':
			return true
		}
		return false
	})
	seen := map[string]struct{}{}
	tickers := make([]string, 0, len(tokens))
	for _, token := range tokens {
		upper := strings.ToUpper(token)
		if len(upper) < 2 || len(upper) > 5 {
			continue
		}
		if _, ok := n.tickerLex[upper]; !ok {
			continue
		}
		if _, exists := seen[upper]; exists {
			continue
		}
		seen[upper] = struct{}{}
		tickers = append(tickers, upper)
	}

	if len(tickers) == 0 {
		return []string{"SPY"}
	}

	return tickers
}

func deterministicID(link, title string, published time.Time) string {
	sum := link + "|" + title + "|" + published.Format(time.RFC3339)
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(sum)).String()
}

func truncate(input string, limit int) string {
	if len(input) <= limit {
		return input
	}
	return input[:limit]
}

func sentimentLabel(score float64) string {
	switch {
	case score > 0.2:
		return "bullish"
	case score < -0.2:
		return "bearish"
	default:
		return "neutral"
	}
}
