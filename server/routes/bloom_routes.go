package routes

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"server/models"
	"server/repository"

	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// type NewsArticle struct {
// 	Title       string `json:"title"`
// 	Link        string `json:"link"`
// 	Description string `json:"description"`
// }

func fetchBloombergNews() ([]models.NewsArticle, error) {
	var articles []models.NewsArticle

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
			article := &models.NewsArticle{
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

// Handler to Fetch and Save News
func getBloombergNews(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		articles, err := fetchBloombergNews()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch news"})
		}

		// Save to Database (only new articles)
		if err := repository.SaveNewsArticles(db, articles); err != nil {
			log.Println("Failed to save news:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save news"})
		}

		return c.JSON(http.StatusOK, articles)
	}
}

func getSavedBloombergNews(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		articles, err := repository.GetAllNewsArticles(db)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve news"})
		}
		return c.JSON(http.StatusOK, articles)
	}
}

func searchNewsArticles(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		query := c.QueryParam("q") // Parameter
		// fmt.Println("Parameter: ", query)

		// Decode Japanese
		decodedQuery, err := url.QueryUnescape(query)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid query"})
		}

		if decodedQuery == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Search query is required"})
		}

		articles, err := repository.SearchNewsArticles(db, decodedQuery)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to search news"})
		}

		// DEBUG:
		// fmt.Println("Results: ", articles)

		return c.JSON(http.StatusOK, articles)
	}
}

// Register Bloomberg Routes
func RegisterBloombergRoutes(e *echo.Group, db *gorm.DB) {
	e.GET("/bloomberg", getBloombergNews(db))
	e.GET("/bloomberg/saved", getSavedBloombergNews(db))
	e.GET("/bloomberg/search", searchNewsArticles(db))
}
