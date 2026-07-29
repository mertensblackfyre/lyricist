package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
)

func main() {

	ctx := context.Background()
	url := flag.String("url", "none", "a string")
	typer := flag.String("type", "none", "a string")

	flag.Parse()

	if *url == "none" || *typer == "none" {
		fmt.Println("Please provide url")
		return
	}

	scraper := NewScraper()
	if *typer == "track" {
		track, err := scraper.ScrapeTrackInfo(ctx, *url)
		if err != nil {
			log.Fatal(err)
		}
		t, err := GetTrack(ctx, &track)
		if err != nil {
			fmt.Println(err)
		}
		query := BuildURL(&t)
		info, err := Search(ctx, query)

		yt_duration := *info.Duration
		deez_duration := float64(t.Duration)

		if math.Abs(yt_duration-deez_duration) > 5 {
			fmt.Printf("duration difference is more than 5, skipping\n")
			return
		}

		DownloadTrack(ctx, *info.WebpageURL)
		/*
			if strings.Contains(strings.ToLower(*info.Title), "cover") {
				fmt.Printf("contains 'cover' keyword ,skipping\n")
			}


		*/
	}
}
