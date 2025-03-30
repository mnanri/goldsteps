package routes

import (
	"fmt"
	"goldsteps/db"
	"goldsteps/models"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
)

// Stock API
type StockData struct {
	Code         string  `json:"code"`
	StockPrice   string  `json:"stock_price"`
	MarketCap    string  `json:"market_cap"`
	IssuedShares string  `json:"issued_shares"`
	PrevClose    string  `json:"prev_close"`
	PriceChange  string  `json:"price_change"`
	StopHigh     bool    `json:"stop_high"`
	AveragePER   float64 `json:"average_per"`
	AveragePBR   float64 `json:"average_pbr"`
}

type Article struct {
	Title  string `json:"title"`
	Link   string `json:"link"`
	Source string `json:"source"`
	Date   string `json:"date"`
}

func stockDailyValue(code string) (StockData, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("minkabu.jp"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*minkabu.jp",
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	// Initialize variables
	var stockData StockData
	stockData.Code = code
	stopHigh := false

	// Extract stock information from Minkabu
	c.OnHTML(".stock_price", func(e *colly.HTMLElement) {
		stockData.StockPrice = strings.TrimSpace(e.Text)
		stockData.StockPrice = strings.ReplaceAll(stockData.StockPrice, "\n", "")
		stockData.StockPrice = strings.TrimSpace(strings.Split(stockData.StockPrice, "円")[0])
		stockData.StockPrice = strings.ReplaceAll(stockData.StockPrice, ",", "")
	})

	c.OnHTML("table.md_table tbody tr", func(e *colly.HTMLElement) {
		label := strings.TrimSpace(e.ChildText("th"))
		value := strings.TrimSpace(e.ChildText("td"))

		// Identify the data based on the label
		switch label {
		case "時価総額": // Market capitalization
			value = strings.ReplaceAll(value, "百万円", ",000,000")
			// value = strings.ReplaceAll(value, ",", "")
			stockData.MarketCap = value
		case "発行済株数": // Issued shares
			value = strings.ReplaceAll(value, "千株", ",000")
			// value = strings.ReplaceAll(value, ",", "")
			stockData.IssuedShares = value
		}
	})

	c.OnHTML("table.md_table.theme_light tr.ly_vamd", func(e *colly.HTMLElement) {
		label := strings.TrimSpace(e.ChildText("th"))
		value := strings.TrimSpace(e.ChildText("td"))

		if strings.Contains(label, "前日終値") { // Match "前日終値"
			value = strings.ReplaceAll(value, "円", "")
			value = strings.ReplaceAll(value, ",", "")
			stockData.PrevClose = value
		}
	})

	// Extract price change and check if "STOP高" exists
	c.OnHTML(".md_stockBoard_stockTable", func(e *colly.HTMLElement) {
		stockData.PriceChange = strings.TrimSpace(e.ChildText(".stock_price_diff"))
		if e.ChildText(".hi") == "STOP高" {
			stopHigh = true
		}
	})

	// Variables for financial data
	var perSum, pbrSum float64
	var count float64

	// Extract PER and PBR values
	c.OnHTML("table.md_table tr", func(e *colly.HTMLElement) {
		cells := e.DOM.Find("td").Map(func(i int, s *goquery.Selection) string {
			return strings.TrimSpace(s.Text())
		})

		if len(cells) >= 4 {
			per, err1 := strconv.ParseFloat(strings.ReplaceAll(cells[2], ",", ""), 64)
			pbr, err2 := strconv.ParseFloat(strings.ReplaceAll(cells[3], ",", ""), 64)
			if err1 == nil && err2 == nil {
				perSum += per
				pbrSum += pbr
				count++
			}
		}
	})

	c.OnScraped(func(r *colly.Response) {
		if count > 0 {
			stockData.AveragePER = math.Round((perSum/count)*100) / 100
			stockData.AveragePBR = math.Round((pbrSum/count)*100) / 100
		} else {
			stockData.AveragePER = 0
			stockData.AveragePBR = 0
		}
		stockData.StopHigh = stopHigh
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Failed to crawl %s: %v", r.Request.URL, err)
	})

	// Visit the stock page
	yahoo_url := fmt.Sprintf("https://minkabu.jp/stock/%s", code)
	c.Visit(yahoo_url)
	minkabu_url := fmt.Sprintf("https://minkabu.jp/stock/%s/daily_valuation", code)
	c.Visit(minkabu_url)

	// For async processing
	time.Sleep(5 * time.Second)

	return stockData, nil
}

// Handler for stock daily value
func getStockInfo(c echo.Context) error {
	// For Debugging
	// fmt.Println("Request Body:", c.Request().Body)
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Stock code is required"})
	}

	var stock models.Stock
	var stockDetail models.StockDetail
	if err := db.DB.Where("stock_code = ?", code).First(&stock).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Stock not found"})
	}

	if err := db.DB.Where("stock_code = ?", code).First(&stockDetail).Error; err != nil {
		log.Println("The stock might be vernished from market?, CODE: ", code, err)
	}

	stockData, err := stockDailyValue(code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch stock data"})
	}

	response := map[string]interface{}{
		"stock":       stock,
		"stockDetail": stockDetail,
		"stockData":   stockData,
	}

	// DEBUG
	// log.Println(response)

	return c.JSON(http.StatusOK, response)
}

// Judge if the date is within one year
func isWithinOneYear(dateStr string) (bool, time.Time) {
	layoutFull := "2006/01/02" // Format for full date (YYYY/MM/DD)
	// layoutShort := "01/02"     // Format for month/day (MM/DD)

	// Remove time part (HH:mm) if present
	dateParts := strings.Split(dateStr, " ")
	cleanDate := dateParts[0] // Extract only the date part

	// Get current date
	now := time.Now()
	oneYearAgo := now.AddDate(-1, 0, 0)

	var parsedDate time.Time
	var err error

	// Try parsing with full date format (YYYY/MM/DD)
	parsedDate, err = time.Parse(layoutFull, cleanDate)
	if err == nil {
		return parsedDate.After(oneYearAgo), parsedDate
	}

	// If parsing fails, try short format (MM/DD) and assume it's from this year or last year
	thisYear := now.Year()
	lastYear := now.Year() - 1

	parsedDate, err = time.Parse(layoutFull, fmt.Sprintf("%d/%s", thisYear, cleanDate))
	if err == nil && parsedDate.After(oneYearAgo) {
		return true, parsedDate
	}

	parsedDate, err = time.Parse(layoutFull, fmt.Sprintf("%d/%s", lastYear, cleanDate))
	if err == nil && parsedDate.After(oneYearAgo) {
		return true, parsedDate
	}

	log.Printf("Date parsing failed for: %s", dateStr)
	return false, time.Time{}
}

func stockNews(code string) ([]Article, error) {
	var articles []Article

	baseURL := "https://minkabu.jp/stock/" + code + "/news?page="
	page := 1
	stopCrawling := false

	for {
		if stopCrawling {
			break
		}

		url := fmt.Sprintf("%s%d", baseURL, page)
		fmt.Println("Visiting:", url)

		c := colly.NewCollector(
			colly.AllowedDomains("minkabu.jp"),
		)

		c.Limit(&colly.LimitRule{
			DomainGlob:  "*minkabu.jp",
			Delay:       2 * time.Second,
			RandomDelay: 1 * time.Second,
		})

		pageHasArticles := false

		c.OnHTML(".md_card_ti", func(e *colly.HTMLElement) {
			if strings.Contains(e.Text, "ページが見つかりませんでした") {
				fmt.Println("Page not found. Stopping crawl.")
				stopCrawling = true
			}
		})

		c.OnHTML("li", func(e *colly.HTMLElement) {
			title := strings.TrimSpace(e.ChildText(".title_box a"))
			link := e.ChildAttr(".title_box a", "href")
			source := strings.TrimSpace(e.ChildText(".fcgl"))
			date := strings.TrimSpace(e.ChildText(".flex.items-center"))

			if source == "適時開示" || source == "PR TIMES" {
				isRecent, articleDate := isWithinOneYear(date)
				if !isRecent {
					fmt.Printf("Found old article (%s), stopping crawl.\n", articleDate.Format("2006/01/02"))
					stopCrawling = true
					return
				}

				// Debug
				// fmt.Println("DEBUG: Extracted Title:", title)
				// fmt.Println("DEBUG: Extracted Link:", link)
				// fmt.Println("DEBUG: Extracted Source:", source)
				// fmt.Println("DEBUG: Extracted Date:", date)

				articles = append(articles, Article{
					Title:  title,
					Link:   "https://minkabu.jp" + link,
					Source: source,
					Date:   date,
				})

				pageHasArticles = true
			}
		})

		if !pageHasArticles {
			fmt.Println("No articles found on this page: ", page)
		}

		err := c.Visit(url)
		if err != nil {
			log.Println("Failed to visit:", err)
			break
		}

		page++
	}

	return articles, nil
}

// Handler for stock news
func getStockNews(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Stock code is required"})
	}

	articles, err := stockNews(code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch news"})
	}

	return c.JSON(http.StatusOK, articles)
}

// RegisterStockRoutes registers stock routes
func RegisterStockRoutes(e *echo.Group) {
	e.GET("/stocks/:code", getStockInfo)
	e.GET("/stocks/:code/news", getStockNews)
}
