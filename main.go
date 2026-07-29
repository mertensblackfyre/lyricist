package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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

		GetTrack(ctx, &track)
		//	query := BuildURL(&track)
		//Search(ctx, query)
		/*
			if strings.Contains(strings.ToLower(*info.Title), "cover") {
				fmt.Printf("contains 'cover' keyword ,skipping\n")
			}

			youtubeDuration := *info.Duration
			spotifyDurationSeconds := float64(track.DurationMs) / 1000.0

			fmt.Println(youtubeDuration, spotifyDurationSeconds, track.DurationMs)
			if math.Abs(youtubeDuration-spotifyDurationSeconds) <= 10.0 {
				fmt.Printf("duration difference is more than 5, skipping\n")
			}

			fmt.Println(*info.Title)
		*/
	}
}
