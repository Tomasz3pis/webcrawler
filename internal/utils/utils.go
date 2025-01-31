package utils

import (
	"cmp"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
)

func PrintReport(pages map[string]int, baseURL string) {
	fmt.Print("=================================\n")
	fmt.Printf("REPORT for %s\n", baseURL)
	fmt.Print("=================================\n")

	//TODO sort by value and if value is the same, by key
	type Result struct {
		pg    string
		count int
	}
	var p []Result
	for k, v := range pages {
		p = append(p, Result{
			pg:    k,
			count: v,
		})
	}
	slices.SortFunc(p, func(a, b Result) int {
		return cmp.Or(
			cmp.Compare(b.count, a.count),
			cmp.Compare(a.pg, b.pg),
		)
	})
	for _, v := range p {
		fmt.Printf("Found %v internal links to %s\n", v.count, v.pg)
	}
}

func GetHTML(rawURL string) (string, error) {
	if !shouldCrawlURL(rawURL) {
		return "", fmt.Errorf("Not a html, skipping")
	}
	r, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}
	if r.StatusCode >= 400 {
		return "", fmt.Errorf("Something went wrong")
	}
	ct := r.Header.Get("content-type")
	if !strings.Contains(ct, "text/html") {
		return "", fmt.Errorf("No content for us")
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func shouldCrawlURL(rawCurrentURL string) bool {
	if strings.HasSuffix(rawCurrentURL, ".xml") ||
		strings.HasSuffix(rawCurrentURL, ".png") ||
		strings.HasSuffix(rawCurrentURL, ".css") ||
		strings.HasSuffix(rawCurrentURL, ".js") {
		return false // Skip non-HTML files
	}

	return true // Crawl this URL
}
