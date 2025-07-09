package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"

	"web-scraper.dev/internal/model"
	"web-scraper.dev/internal/repository"
	"web-scraper.dev/internal/tasks"
	l "web-scraper.dev/internal/utils/logger"
)

const (
	scraperUserAgent        = `"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"`
	fmtScrapeKeywordPageUrl = "https://www.bing.com/search?q=%s"
)

type ScrapingResult struct {
	AdCount     int
	LinkCount   int
	HTMLContent string
}

type ScrapeWorker struct {
	srv    *asynq.Server
	db     *repository.Db
	logger *l.Logger
}

func NewScrapeWorker(redisOpt asynq.RedisClientOpt, db *gorm.DB, logger *l.Logger) *ScrapeWorker {
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"default": 10,
			},
		},
	)

	return &ScrapeWorker{
		srv:    srv,
		db:     repository.New(db),
		logger: logger,
	}
}

func (w *ScrapeWorker) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeScrapeKeyword, w.HandleSearchScrapeTask)
	return w.srv.Run(mux)
}

func (w *ScrapeWorker) Stop() error {
	w.srv.Shutdown()
	return nil
}

func (w *ScrapeWorker) HandleSearchScrapeTask(ctx context.Context, t *asynq.Task) error {
	var p tasks.ScrapeKeywordPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		w.logger.Error().Err(err).Msg("failed to unmarshal payload")
		return err
	}

	keywordID := p.KeywordID

	if err := w.db.Model(&model.Keyword{}).Where("id = ?", keywordID).Update("status", "processing").Error; err != nil {
		w.logger.Error().Err(err).Msg("failed to update keyword status")
		return err
	}

	var keyword model.Keyword
	if err := w.db.First(&keyword, keywordID).Error; err != nil {
		w.logger.Error().Err(err).Msg("failed to find keyword")
		return err
	}

	result, err := w.scrapeKeyword(keyword.Keyword)
	if err != nil {
		w.logger.Error().Err(err).Msg("failed to scrape keyword")
		w.updateKeywordError(keywordID, err.Error())
		return err
	}

	updates := map[string]interface{}{
		"status":       "completed",
		"ad_count":     result.AdCount,
		"link_count":   result.LinkCount,
		"html_content": html.EscapeString(result.HTMLContent),
	}

	if err := w.db.Model(&model.Keyword{}).Where("id = ?", keywordID).Updates(updates).Error; err != nil {
		w.logger.Error().Err(err).Msg("failed to update keyword results")
		return err
	}

	w.logger.Info().Msgf("Successfully processed KeywordID: %d", keywordID)
	return nil
}

func (w *ScrapeWorker) updateKeywordError(keywordID int64, errorMsg string) {
	updates := map[string]interface{}{
		"status":        "failed",
		"error_message": errorMsg,
	}

	if err := w.db.Model(&model.Keyword{}).Where("id = ?", keywordID).Updates(updates).Error; err != nil {
		w.logger.Error().Err(err).Msg("failed to update keyword error")
	}
}

func (w *ScrapeWorker) scrapeKeyword(keyword string) (*ScrapingResult, error) {
	c := colly.NewCollector(
		colly.UserAgent(scraperUserAgent),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*bing.*",
		Parallelism: 1,
		Delay:       tasks.ScrapeKeywordDelayInSeconds * time.Second,
	})

	c.SetRequestTimeout(30 * time.Second)

	var result ScrapingResult
	var htmlContent strings.Builder

	// Count Bing ads
	c.OnHTML(".b_ad, .b_adurl, .sb_add, [data-bm], .b_adSlug", func(e *colly.HTMLElement) {
		result.AdCount++
	})

	// Count all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if href != "" && !strings.HasPrefix(href, "#") {
			result.LinkCount++
		}
	})

	// Capture HTML content
	c.OnResponse(func(r *colly.Response) {
		htmlContent.Write(r.Body)
	})

	// Construct Bing search URL
	searchURL := fmt.Sprintf(fmtScrapeKeywordPageUrl, url.QueryEscape(keyword))

	// Visit the URL
	err := c.Visit(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to visit URL: %v", err)
	}

	c.Wait()

	// Set HTML content
	htmlStr := htmlContent.String()
	if len(htmlStr) > 1024*1024 {
		htmlStr = htmlStr[:1024*1024] + "... [truncated]"
	}
	result.HTMLContent = htmlStr

	return &result, nil
}
