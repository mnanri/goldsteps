package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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
		// colly.CacheDir("./colly_cache"),             // Use cache to prevent duplicate visits
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
				// fmt.Printf("Link found: %s\nTitle: %s\n", absoluteURL, title)

				// Visit the link if it hasn't been visited already
				visited, err := c.HasVisited(absoluteURL)
				if err != nil {
					log.Println("Error checking visit status:", err)
					return
				}
				if !visited {
					fmt.Printf("Visiting: %s\n", absoluteURL)
					fmt.Printf("Article found: %s\n", title)
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

func minkabuListedStocksFundamental(filename string) {
	c := colly.NewCollector(
		colly.AllowedDomains("minkabu.jp"),
		colly.CacheDir("./colly_cache"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*minkabu.jp",
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	// Write headers to CSV file
	headers := []string{"銘柄コード", "上場区分", "銘柄", "株価", "社名", "英文社名", "業種", "代表者", "決算", "資本金", "住所", "電話番号(IR)", "上場市場", "上場年月日", "単元株数"}
	var record []string

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
		record = []string{stockCode, listingSection, stockName, stockPrice, "", "", "", "", "", "", "", "", "", "", ""}

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
				record[4] = value
			case "英文社名":
				fmt.Println("英文社名:", value)
				record[5] = value
			case "業種":
				fmt.Println("業種:", value)
				record[6] = value
			case "代表者":
				fmt.Println("代表者:", value)
				record[7] = value
			case "決算":
				fmt.Println("決算:", value)
				record[8] = value
			case "資本金":
				fmt.Println("資本金:", value)
				record[9] = value
			case "住所":
				fmt.Println("住所:", value)
				record[10] = value
			case "電話番号(IR)":
				fmt.Println("電話番号(IR):", value)
				record[11] = value
			case "上場市場":
				fmt.Println("上場市場:", value)
				record[12] = value
			case "上場年月日":
				fmt.Println("上場年月日:", value)
				record[13] = value
			case "単元株数":
				fmt.Println("単元株数:", value)
				record[14] = value
			}
		})
		if record[12] != "" || record[13] != "" || record[14] != "" {
			writeToCSV(filename, headers, record)
			fmt.Println("--------------------------------------------------")
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Failed to crawl %s: %v", r.Request.URL, err)
	})

	for code := 1300; code <= 9999; code++ {
		url := fmt.Sprintf("https://minkabu.jp/stock/%d", code)
		c.Visit(url)
	}
}

func yahooFinanceStockProfile(filename string, stock_list [][]string) {
	// Colly instance
	c := colly.NewCollector(
		colly.AllowedDomains("finance.yahoo.co.jp"),
	)

	// Limit the request rate
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*yahoo.co.jp",
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	// Write the company name to a CSV file
	headers := []string{"特色", "連結事業", "従業員数（単独）", "従業員数（連結）", "平均年齢", "平均年収"}
	record := make([]string, len(headers))

	// Extract information based on the corresponding table headers
	c.OnHTML("table.CompanyInformationDetail__table__BIq9", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, row *colly.HTMLElement) {
			header := row.ChildText("th")
			value := row.ChildText("td")

			// Clean up the extracted text
			value = strings.TrimSpace(value)
			value = strings.ReplaceAll(value, "【特色】", "")
			value = strings.ReplaceAll(value, "【連結事業】", "")
			if strings.Contains(value, "人") {
				value = strings.ReplaceAll(value, "人", "")
				value = strings.ReplaceAll(value, ",", "")
			}
			if strings.Contains(value, "歳") {
				value = strings.ReplaceAll(value, "歳", "")
			}
			if strings.Contains(value, "円") {
				value = strings.ReplaceAll(value, "円", "")
				value = strings.ReplaceAll(value, ",", "")
				value = strings.ReplaceAll(value, "千", "000")
				value = strings.ReplaceAll(value, "百万", "000000")
			}

			switch header {
			case "特色":
				record[0] = value
				fmt.Println("特色:", value)
			case "連結事業":
				record[1] = value
				fmt.Println("連結事業:", value)
			case "従業員数（単独）":
				record[2] = value
				fmt.Println("従業員数（単独）:", value)
			case "従業員数（連結）":
				record[3] = value
				fmt.Println("従業員数（連結）:", value)
			case "平均年齢":
				record[4] = value
				fmt.Println("平均年齢:", value)
			case "平均年収":
				record[5] = value
				fmt.Println("平均年収:", value)
			}
		})

		// Write the extracted data to a CSV file
		writeToCSV(filename, headers, record)
	})

	// Error handling
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Failed to crawl %s: %v", r.Request.URL, err)
	})

	for _, stock := range stock_list {
		code := stock[0]
		name := stock[2]
		fmt.Printf("Code: %s\tName: %s\n", code, name)
		url := fmt.Sprintf("https://finance.yahoo.co.jp/quote/%s.T/profile", code)
		c.Visit(url)
		fmt.Println("--------------------------------------------------")
	}
}

func writeToCSV(filename string, headers []string, record []string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if fileInfo, _ := file.Stat(); fileInfo.Size() == 0 {
		writer.Write(headers)
	}

	writer.Write(record)
}

func readFromCSV(filename string) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %v", err)
	}

	return records
}

func stockDailyValue(code string) bool {
	c := colly.NewCollector(
		colly.AllowedDomains("minkabu.jp"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*minkabu.jp",
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	// Variables for stock information
	var stockPrice, marketCap, issuedShares, prevClose, priceChange string
	stopHigh := false

	// Extract stock information from Minkabu
	c.OnHTML(".stock_price", func(e *colly.HTMLElement) {
		stockPrice = strings.TrimSpace(e.Text)
		stockPrice = strings.ReplaceAll(stockPrice, "\n", "")
		stockPrice = strings.TrimSpace(strings.Split(stockPrice, "円")[0])
		stockPrice = strings.ReplaceAll(stockPrice, ",", "")
	})

	c.OnHTML("table.md_table tbody tr", func(e *colly.HTMLElement) {
		label := strings.TrimSpace(e.ChildText("th"))
		value := strings.TrimSpace(e.ChildText("td"))

		// Identify the data based on the label
		switch label {
		case "時価総額": // Market capitalization
			value = strings.ReplaceAll(value, "百万円", "000000")
			value = strings.ReplaceAll(value, ",", "")
			marketCap = value
		case "発行済株数": // Issued shares
			value = strings.ReplaceAll(value, "千株", "000")
			value = strings.ReplaceAll(value, ",", "")
			issuedShares = value
		}
	})

	c.OnHTML("table.md_table.theme_light tr.ly_vamd", func(e *colly.HTMLElement) {
		label := strings.TrimSpace(e.ChildText("th"))
		value := strings.TrimSpace(e.ChildText("td"))

		if strings.Contains(label, "前日終値") { // Match "前日終値"
			value = strings.ReplaceAll(value, "円", "")
			value = strings.ReplaceAll(value, ",", "")
			prevClose = value
		}
	})

	// Extract price change and check if "STOP高" exists
	c.OnHTML(".md_stockBoard_stockTable", func(e *colly.HTMLElement) {
		priceChange = strings.TrimSpace(e.ChildText(".stock_price_diff"))
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
			fmt.Println("Stock Information for Code:", code)
			fmt.Printf("Stock Price: %s\n", stockPrice)
			fmt.Printf("Market Capitalization: %s\n", marketCap)
			fmt.Printf("Issued Shares: %s\n", issuedShares)
			fmt.Printf("Previous Closing Price: %s\n", prevClose)
			fmt.Printf("Price Change: %s\n", priceChange)
			if stopHigh {
				fmt.Println("STOP高: Yes")
			} else {
				fmt.Println("STOP高: No")
			}
			fmt.Printf("Ave. PER: %.2f\n", perSum/count)
			fmt.Printf("Ave. PBR: %.2f\n", pbrSum/count)
		} else {
			fmt.Println("Loading financial data...")
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Failed to crawl %s: %v", r.Request.URL, err)
	})

	// Visit the stock page
	yahoo_url := fmt.Sprintf("https://minkabu.jp/stock/%s", code)
	c.Visit(yahoo_url)
	minkabu_url := fmt.Sprintf("https://minkabu.jp/stock/%s/daily_valuation", code)
	c.Visit(minkabu_url)

	return stopHigh
}

// Function to check if the date is within this year only
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

type Article struct {
	Title  string `json:"title"`
	Link   string `json:"link"`
	Source string `json:"source"`
	Date   string `json:"date"`
}

func stockNews(code string) {
	var articles []Article // List to store articles

	baseURL := "https://minkabu.jp/stock/" + code + "/news?page="
	page := 1
	stopCrawling := false

	for {
		if stopCrawling {
			break
		}

		// Ex. https://minkabu.jp/stock/4385/news?page=1
		url := fmt.Sprintf("%s%d", baseURL, page)
		fmt.Println("Visiting:", url)

		c := colly.NewCollector(
			colly.AllowedDomains("minkabu.jp"),
			colly.CacheDir("./colly_cache"),
		)

		c.Limit(&colly.LimitRule{
			DomainGlob:  "*minkabu.jp",
			Delay:       2 * time.Second,
			RandomDelay: 1 * time.Second,
		})

		pageHasArticles := false // Track if this page contains valid articles

		// Detect if the page does not exist
		c.OnHTML(".md_card_ti", func(e *colly.HTMLElement) {
			if strings.Contains(e.Text, "ページが見つかりませんでした") {
				fmt.Println("Page not found. Stopping crawl.")
				stopCrawling = true
			}
		})

		// Extract news items
		c.OnHTML("li", func(e *colly.HTMLElement) {
			title := strings.TrimSpace(e.ChildText(".title_box a"))
			link := e.ChildAttr(".title_box a", "href")
			source := strings.TrimSpace(e.ChildText(".fcgl"))
			date := strings.TrimSpace(e.ChildText(".flex.items-center"))

			// Only process "適時開示" or "PR TIMES" articles
			if source == "適時開示" || source == "PR TIMES" {
				isRecent, articleDate := isWithinOneYear(date)

				if !isRecent {
					fmt.Printf("Found old article (%s), stopping crawl.\n", articleDate.Format("2006/01/02"))
					stopCrawling = true
					return
				}

				fmt.Println("Title:", title)
				fmt.Println("Link: https://minkabu.jp" + link)
				fmt.Println("Source:", source)
				fmt.Println("Date:", date)
				fmt.Println("----------------------")

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
			fmt.Println("No articles found on this page.")
		}

		// Visit the page
		err := c.Visit(url)
		if err != nil {
			log.Fatal(err)
		}

		// Move to the next page, even if the current page has no articles
		page++
	}

	// Convert to JSON and write to file
	writeJSON(articles, code)
}

// Function to write articles to a JSON file
func writeJSON(articles []Article, code string) {
	today := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("articles_%s_%s.json", code, today)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Failed to create JSON file:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print JSON
	err = encoder.Encode(articles)
	if err != nil {
		log.Fatal("Failed to write JSON:", err)
	}

	fmt.Println("Articles saved to articles.json")
}

func main() {
	// Get the top news from Bloomberg
	// bloomTopNews()
	bloomTopNewsDescription()

	// Get the listed stocks and their fundamental information
	// minkabu_stock_fundamental_filename := "stock_fundamental_202502.csv"
	// minkabuListedStocks()
	// minkabuListedStocksFundamental(minkabu_stock_fundamental_filename)

	// Display the data
	// records := readFromCSV(minkabu_stock_fundamental_filename)
	// for _, record := range records {
	// 	fmt.Printf("Code: %s\tName: %s\tPrice: %s\n", record[0], record[2], record[3])
	// }

	// Get the stock profile from Yahoo Finance
	// yahoo_finance_stock_profile_filename := "stock_profile.csv"
	// yahooFinanceStockProfile(yahoo_finance_stock_profile_filename, records)

	// Get the daily stock value
	// mercari_code := "4385"
	// stockDailyValue(mercari_code)

	// stockNews(mercari_code)
}
