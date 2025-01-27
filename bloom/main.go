package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Colly Instance
	c := colly.NewCollector(
		colly.AllowedDomains("www.bloomberg.co.jp"), // Restrict to specific domain
	)

	// Limit the request rate
	c.Limit(&colly.LimitRule{
		DomainGlob: "*bloomberg.co.jp",
		Delay:      2 * time.Second, // 2秒の遅延を設定
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
