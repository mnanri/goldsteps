package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	// bloomTopLinks()
	// fmt.Println("--------------------------------------------------")
	bloomTopLinksDescription()
}

func bloomTopLinks() {
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

func bloomTopLinksDescription() {
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
