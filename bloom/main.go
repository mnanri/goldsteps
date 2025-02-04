package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func bloomTopNews() {
	// Colly Instance
	c := colly.NewCollector(
		colly.AllowedDomains("www.bloomberg.co.jp"), // Restrict to specific domain
	)

	// Limit the request rate
	c.Limit(&colly.LimitRule{
		DomainGlob: "*bloomberg.co.jp",
		Delay:      2 * time.Second, // Set a delay between requests
	})

	// Extract specific elements
	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println("Page Title:", e.Text)
	})

	// Extract and print all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")                      // Get the href attribute
		absoluteURL := e.Request.AbsoluteURL(link)  // Convert to absolute URL
		title := e.Text                             // Get the link text (title)
		title = strings.TrimSpace(title)            // Remove leading and trailing white spaces
		title = strings.ReplaceAll(title, "\n", "") // Remove newlines
		title = strings.ReplaceAll(title, "\t", "") // Remove tabs
		if title != "" && strings.Contains(absoluteURL, "https://www.bloomberg.co.jp/news/articles") {
			fmt.Printf("Link found: %s\nTitle: %s\n", absoluteURL, title)
		}
	})

	// Error Handling
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request failed:", err)
	})

	// Visiting
	err := c.Visit("https://www.bloomberg.co.jp/")
	if err != nil {
		log.Fatal(err)
	}
}

// Struct to parse JSON-LD data
type NewsArticle struct {
	Description string `json:"description"`
}

func bloomTopNewsDescription() {
	// Create a new Colly collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.bloomberg.co.jp"), // Restrict crawling to specific domain
		colly.CacheDir("./colly_cache"),             // Use cache to prevent duplicate visits
	)

	// Limit request rate to avoid overloading the server
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*bloomberg.co.jp",
		Delay:       2 * time.Second, // Delay between requests
		RandomDelay: 1 * time.Second, // Random delay for load balancing
	})

	// Extract and print the page title
	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println("Page Title:", e.Text)
	})

	// Extract and visit article links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Request.URL.String() == "https://www.bloomberg.co.jp/" {
			link := e.Attr("href")
			absoluteURL := e.Request.AbsoluteURL(link) // Convert to absolute URL
			title := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(e.Text, "\n", ""), "\t", ""))

			if title != "" && strings.Contains(absoluteURL, "https://www.bloomberg.co.jp/news/articles") {
				fmt.Printf("Link found: %s\nTitle: %s\n", absoluteURL, title)

				// Visit the link if it hasn't been visited already
				visited, err := c.HasVisited(absoluteURL)
				if err != nil {
					log.Println("Error checking visit status:", err)
					return
				}
				if !visited {
					err := c.Visit(absoluteURL)
					if err != nil {
						log.Println("Visit failed:", err)
					}
				}
			}
		}
	})

	// Extract description from JSON-LD script
	c.OnHTML(`script[type="application/ld+json"]`, func(e *colly.HTMLElement) {
		var article NewsArticle
		err := json.Unmarshal([]byte(e.Text), &article)
		if err == nil && article.Description != "" {
			fmt.Println("Description:", article.Description)
			fmt.Println("--------------------------------------------------")
		}
	})

	// Handle errors during the request
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with response: %v\nError: %v", r.Request.URL, r, err)
	})

	// Start crawling from the homepage
	err := c.Visit("https://www.bloomberg.co.jp/")
	if err != nil {
		log.Fatal(err)
	}
}

// Get listed stocks
func minkabuListedStocks() {
	c := colly.NewCollector(
		colly.AllowedDomains("minkabu.jp"),
		colly.CacheDir("./colly_cache"), // Prevent duplicate visits
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*minkabu.jp",
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	// Extract stock information
	c.OnHTML("div.md_stockBoard", func(e *colly.HTMLElement) {
		// Stock code from URL
		stockCode := strings.Split(e.Request.URL.Path, "/")[2]

		// Listing section (e.g., 東証プライム)
		listingSection := e.ChildText("div.stock_label")
		listingSection = strings.TrimSpace(strings.Split(listingSection, "  ")[1]) // Use full-width space

		// Stock name (e.g., メルカリ)
		stockName := e.ChildText("h2 span.md_stockBoard_stockName")

		// Stock price (e.g., 1,862円)
		stockPrice := e.ChildText("div.stock_price")
		stockPrice = strings.ReplaceAll(stockPrice, "\n", "")
		stockPrice = strings.TrimSpace(stockPrice)
		stockPrice = strings.TrimSpace(strings.Split(stockPrice, "円")[0])
		stockPrice = strings.ReplaceAll(stockPrice, ",", "")
		// Output
		fmt.Printf("銘柄コード: %s\n上場区分: %s\n銘柄: %s\n株価: %s\n", stockCode, listingSection, stockName, stockPrice)
		fmt.Println("--------------------------------------------------")
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Failed to crawl %s: %v", r.Request.URL, err)
	})

	// Crawl stock pages from 1000 to 9999
	for code := 4350; code <= 4400; code += 5 {
		url := fmt.Sprintf("https://minkabu.jp/stock/%d", code)
		c.Visit(url)
	}
}

func minkabuListedStocksFundamental() {
	c := colly.NewCollector(
		colly.AllowedDomains("minkabu.jp"),
		colly.CacheDir("./colly_cache"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*minkabu.jp",
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	// Extract stock information from the main page
	c.OnHTML("div.md_stockBoard", func(e *colly.HTMLElement) {
		stockCode := strings.Split(e.Request.URL.Path, "/")[2]
		listingSection := e.ChildText("div.stock_label")
		listingSection = strings.TrimSpace(strings.Split(listingSection, "  ")[1])
		stockName := e.ChildText("h2 span.md_stockBoard_stockName")
		stockPrice := e.ChildText("div.stock_price")
		stockPrice = strings.ReplaceAll(stockPrice, "\n", "")
		stockPrice = strings.TrimSpace(stockPrice)
		stockPrice = strings.TrimSpace(strings.Split(stockPrice, "円")[0])
		stockPrice = strings.ReplaceAll(stockPrice, ",", "")
		// Output
		fmt.Printf("銘柄コード: %s\n上場区分: %s\n銘柄: %s\n株価: %s\n", stockCode, listingSection, stockName, stockPrice)

		// Visit the fundamental page for more details
		fundamentalURL := fmt.Sprintf("https://minkabu.jp/stock/%s/fundamental", stockCode)
		e.Request.Visit(fundamentalURL)
	})

	// Extract company fundamental information
	c.OnHTML("dl.md_dataList", func(e *colly.HTMLElement) {
		e.ForEach("dt", func(_ int, dt *colly.HTMLElement) {
			label := dt.Text
			value := dt.DOM.Next().Text()
			value = strings.TrimSpace(value)

			switch label {
			case "社名":
				fmt.Println("社名:", value)
			case "英文社名":
				fmt.Println("英文社名:", value)
			case "業種":
				fmt.Println("業種:", value)
			case "代表者":
				fmt.Println("代表者:", value)
			case "決算":
				fmt.Println("決算:", value)
			case "資本金":
				fmt.Println("資本金:", value)
			case "住所":
				fmt.Println("住所:", value)
			case "電話番号(IR)":
				fmt.Println("電話番号(IR):", value)
			case "上場市場":
				fmt.Println("上場市場:", value)
			case "上場年月日":
				fmt.Println("上場年月日:", value)
			case "単元株数":
				fmt.Println("単元株数:", value)
			}
		})
		fmt.Println("--------------------------------------------------")
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Failed to crawl %s: %v", r.Request.URL, err)
	})

	for code := 4000; code <= 4400; code += 5 {
		url := fmt.Sprintf("https://minkabu.jp/stock/%d", code)
		c.Visit(url)
	}
}

func main() {
	// bloomTopNews()
	// bloomTopNewsDescription()

	// minkabuListedStocks()
	minkabuListedStocksFundamental()
}
