package services

import (
	"context"
	"go-scratch/generated"
	"go-scratch/internal/repository"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

type CrawlerRepository interface {
	InsertOne(crawlResult repository.CrawlResult) error
	FindAll() ([]repository.CrawlResult, error)
}

type CrawlerService struct {
	repo CrawlerRepository
}

func NewCrawlerService(repo CrawlerRepository) *CrawlerService {
	return &CrawlerService{
		repo: repo,
	}
}

func (s *CrawlerService) Crawl(ctx context.Context, url string, depth int) (generated.CrawlResponse, error) {
	// Initialize response
	response := generated.CrawlResponse{
		Url:       url,
		CrawledAt: time.Now(),
		Links:     make([]string, 0),
	}

	// Initialize collector
	c := colly.NewCollector(
		colly.MaxDepth(depth),
	)

	var mu sync.Mutex

	// Extract title
	c.OnHTML("title", func(e *colly.HTMLElement) {
		mu.Lock()
		response.Title = e.Text
		mu.Unlock()
	})

	// Extract links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		mu.Lock()
		response.Links = append(response.Links, link)
		mu.Unlock()
	})

	// Get page content
	c.OnResponse(func(r *colly.Response) {
		mu.Lock()
		response.Content = string(r.Body)
		mu.Unlock()
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		mu.Lock()
		response.Content = err.Error()
		mu.Unlock()
	})

	// Set rate limiting
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 1 * time.Second,
	})

	// Start crawling
	err := c.Visit(url)
	if err != nil {
		return response, err
	}

	// Wait for all crawls to complete
	c.Wait()

	dataDocs := repository.CrawlResult{
		URL:       response.Url,
		Title:     response.Title,
		Content:   response.Content,
		Links:     response.Links,
		CrawledAt: response.CrawledAt,
	}

	if err := s.repo.InsertOne(dataDocs); err != nil {
		return response, err
	}

	return response, nil
}
