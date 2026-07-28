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
		track, err := scraper.GetTrack(ctx, *url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Track: %s by %s\n", track.Title, track.Artists)
	}
}
