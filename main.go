package main

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"webcrawler/internal/utils"
)

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		log.Print("not enough arguments provided\n")
		log.Print("usage: crawler <baseURL> <maxConcurrency> <maxPages>\n")
		os.Exit(1)
	}
	if len(args) > 3 {
		log.Print("too many arguments provided\n")
		os.Exit(1)
	}
	rawURL := args[0]
	mc, err1 := strconv.Atoi(args[1])
	mp, err2 := strconv.Atoi(args[2])
	if err1 != nil || err2 != nil {
		log.Panicf("Provided arguments are not a numbers. %s\n%s\n", err1, err2)
	}
	cfg, err := configure(rawURL, mc, mp)
	if err != nil {
		log.Printf("Error - configure: %v", err)
	}
	log.Printf("starting crawl of: %s\n", rawURL)
	cfg.wg.Add(1)
	go cfg.crawlPage(rawURL)
	cfg.wg.Wait()
	utils.PrintReport(cfg.pages, cfg.baseURL.String())
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()
	if cfg.pagesLen() >= cfg.maxPages {
		return
	}
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		log.Printf("Failed to parse url: %s", err)
		return
	}
	if !strings.HasSuffix(currentURL.Hostname(), cfg.baseURL.Hostname()) {
		return
	}
	normalizedCurrnetURL, err := utils.NormalizeURL(rawCurrentURL)
	if err != nil {
		log.Panicf("Failed to normalize url: %s", err)
	}
	if !cfg.addPageVisit(normalizedCurrnetURL) {
		return
	}
	log.Printf("Crawling %s\n", rawCurrentURL)
	body, err := utils.GetHTML(rawCurrentURL)
	if err != nil {
		log.Printf("Failed to fetch html from: %s: %s", rawCurrentURL, err)
		return
	}
	urls, err := utils.GetURLsFromHTML(body, *cfg.baseURL)
	if err != nil {
		log.Fatalf("Failed to get urls form html: %s", err)
	}
	for _, v := range urls {
		cfg.wg.Add(1)
		go cfg.crawlPage(v)
	}
}
