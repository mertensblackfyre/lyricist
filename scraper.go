package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func ScrapeTrackInfo(ctx context.Context, url string) (TrackScrapeInfo, error) {

	body, err := SendRequest(ctx, url,
		http.Client{
			Timeout: 15 * time.Second,
		})

	if err != nil {
		return TrackScrapeInfo{}, err
	}

	final, err := ParseRecordingJSON(string(body))
	if err != nil {
		return TrackScrapeInfo{}, fmt.Errorf("extracting JSON-LD: %w", err)
	}

	art := ParseArtists(final[1])

	return TrackScrapeInfo{
		Title:   final[0],
		Artists: art,
	}, nil
}

func ScarapePlaylistTracks(ctx context.Context, p_url string) ([]string, error) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/usr/sbin/brave"), // change to your path
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.NoSandbox,
	)
	alloc_ctx, cancel_alloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel_alloc()

	browser_ctx, cancel_browser := chromedp.NewContext(alloc_ctx)
	defer cancel_browser()

	load_ctx, cancel_load := context.WithTimeout(browser_ctx, 30*time.Second)
	defer cancel_load()

	var links []string
	seen := make(map[string]bool)

	err := chromedp.Run(load_ctx,
		chromedp.Navigate(p_url),
		chromedp.WaitVisible(`div[role="row"]`, chromedp.ByQuery),
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('a[href*="/track/"]'))
				.map(a => a.href)
		`, &links),
	)
	if err != nil {
		return nil, fmt.Errorf("headless playlist extraction: %w", err)
	}

	var finalURLs []string
	for _, href := range links {
		u, err := url.Parse(href)
		if err != nil {
			continue
		}
		clean := "https://open.spotify.com" + u.Path
		if !seen[clean] && strings.Contains(clean, "/track/") {
			seen[clean] = true
			finalURLs = append(finalURLs, clean)
		}
	}

	if len(finalURLs) == 0 {
		return nil, fmt.Errorf("no track URLs found on playlist page")
	}
	fmt.Println(finalURLs)
	return finalURLs, nil
}
