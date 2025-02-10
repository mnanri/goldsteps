package routes

import (
	"fmt"
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

	stockData, err := stockDailyValue(code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch stock data"})
	}

	return c.JSON(http.StatusOK, stockData)
}

// RegisterStockRoutes registers stock routes
func RegisterStockRoutes(e *echo.Group) {
	e.GET("/stocks/:code", getStockInfo)
}
