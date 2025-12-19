package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/llan0/link-fetcher/internal/config"
	"github.com/llan0/link-fetcher/internal/fetcher"
	"github.com/llan0/link-fetcher/internal/writer"
)

const outputDir = "output"

func main() {
	runYoutube := flag.Bool("youtube", false, "Run YouTube scraper")
	runTiktok := flag.Bool("tiktok", false, "Run TikTok scraper")
	runInsta := flag.Bool("instagram", false, "Run Instagram scraper")
	runAll := flag.Bool("all", false, "Run ALL scrapers")
	flag.Parse()

	if *runAll {
		*runYoutube = true
		*runTiktok = true
		*runInsta = true
	}

	if !*runYoutube && !*runTiktok && !*runInsta {
		fmt.Println("Usage: make -tiktok (or -youtube, -instagram, -all)")
		return
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Printf("Fetching videos from %s to %s...\n", cfg.StartDateStr, cfg.EndDateStr)

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// YouTube
	if *runYoutube {
		wg.Add(1)
		filename := filepath.Join(outputDir, "youtube.txt")

		go runScraper(ctx, &wg, filename, &fetcher.YoutubeFetcher{
			APIKey:    cfg.APIKeys.Youtube,
			ChannelID: cfg.Accounts.YoutubeChannelID,
			Start:     cfg.StartDate,
			End:       cfg.EndDate,
		})
	}

	// TikTok
	if *runTiktok {
		wg.Add(1)
		filename := filepath.Join(outputDir, "tiktok.txt")

		go runScraper(ctx, &wg, filename, &fetcher.TikTokFetcher{
			RapidAPIKey:  cfg.APIKeys.RapidAPIKey,
			RapidAPIHost: cfg.APIKeys.TikTokHost,
			Username:     cfg.Accounts.TikTokUsername,
			Start:        cfg.StartDate,
			End:          cfg.EndDate,
		})
	}

	// Instagram
	if *runInsta {
		wg.Add(1)
		filename := filepath.Join(outputDir, "instagram.txt")

		go runScraper(ctx, &wg, filename, &fetcher.InstagramFetcher{
			RapidAPIKey:  cfg.APIKeys.RapidAPIKey,
			RapidAPIHost: cfg.APIKeys.InstagramHost,
			Username:     cfg.Accounts.InstagramUsername,
			Start:        cfg.StartDate,
			End:          cfg.EndDate,
		})
	}

	wg.Wait()
	fmt.Println("Done!")
}

func runScraper(ctx context.Context, wg *sync.WaitGroup, filename string, p fetcher.Platform) {
	defer wg.Done()
	fmt.Printf("[%s] Starting...\n", p.Name())

	links, err := p.FetchLinks(ctx)
	if err != nil {
		log.Printf("[%s] Error: %v", p.Name(), err)
		return
	}

	if len(links) == 0 {
		fmt.Printf("[%s] No videos found in this date range.\n", p.Name())
		return
	}

	err = writer.WriteLinks(filename, links)
	if err != nil {
		log.Printf("[%s] Error writing output file: %v", p.Name(), err)
		return
	}

	fmt.Printf("[%s] Success. Saved %d links to %s\n", p.Name(), len(links), filename)
}
