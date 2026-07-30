package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	ctx := context.Background()
	allowedTypes := []string{"track", "album", "playlist"}

	urlFlag := flag.String("url", "", "Spotify URL to process")
	typeFlag := flag.String("type", "", fmt.Sprintf("Type of resource (%s)", strings.Join(allowedTypes, ", ")))
	outputDir := flag.String("output", "output", "Output directory for downloaded files")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s -url <spotify_url> -type <type> [-output <dir>]\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nExample:\n  %s -url https://open.spotify.com/track/xxxx -type track -output ./music\n", os.Args[0])
	}

	flag.Parse()

	if *urlFlag == "" || *typeFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	valid := false
	lowerType := strings.ToLower(*typeFlag)
	for _, t := range allowedTypes {
		if lowerType == t {
			valid = true
			*typeFlag = t
			break
		}
	}
	if !valid {
		fmt.Fprintf(os.Stderr, "Invalid type: %q. Allowed values: %s\n", *typeFlag, strings.Join(allowedTypes, ", "))
		flag.Usage()
		os.Exit(1)
	}

	Create(*outputDir)

	if *typeFlag == "track" {
		err := HandleTrack(ctx, *urlFlag, *outputDir)
		if err != nil {
			log.Println(err)
		}
	}
	Clean()
}
