package routes

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
)

type NewsArticle struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
}

func fetchBloombergNews() ([]NewsArticle, error) {
	var articles []NewsArticle

	// Colly Instance
	c := colly.NewCollector(
		colly.AllowedDomains("www.bloomberg.co.jp"), // Restrict to bloomberg.co.jp
	)

	// Limit the rate of requests
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*bloomberg.co.jp",
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	// Extract article titles and links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteURL := e.Request.AbsoluteURL(link) // Convert relative URL to absolute URL
		title := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(e.Text, "\n", ""), "\t", ""))

		// Filter out empty titles and only include articles
		if title != "" && strings.Contains(absoluteURL, "https://www.bloomberg.co.jp/news/articles") {
			article := &NewsArticle{
				Title: title,
				Link:  absoluteURL,
			}
			articles = append(articles, *article)
		}
	})

	// Error handling
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s, Error: %v", r.Request.URL, err)
	})

	// Start the crawl
	err := c.Visit("https://www.bloomberg.co.jp/")
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// Handler for fetching Bloomberg news
func getBloombergNews(c echo.Context) error {
	articles, err := fetchBloombergNews()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch news"})
	}
	return c.JSON(http.StatusOK, articles)
}

// RegisterBloombergRoutes registers routes for fetching Bloomberg news
func RegisterBloombergRoutes(e *echo.Group) {
	e.GET("/bloomberg", getBloombergNews)
}
