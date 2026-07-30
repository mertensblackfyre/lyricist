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
	allowed_types := []string{"track", "album", "playlist"}

	url_flag := flag.String("url", "", "Spotify URL to process")
	type_flag := flag.String("type", "", fmt.Sprintf("Type of resource (%s)", strings.Join(allowed_types, ", ")))
	output_dir := flag.String("output", "output", "Output directory for downloaded files")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s -url <spotify_url> -type <type> [-output <dir>]\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nExample:\n  %s -url https://open.spotify.com/track/xxxx -type track -output ./music\n", os.Args[0])
	}

	flag.Parse()

	if *url_flag == "" || *type_flag == "" {
		flag.Usage()
		os.Exit(1)
	}

	valid := false
	lower_type := strings.ToLower(*type_flag)
	for _, t := range allowed_types {
		if lower_type == t {
			valid = true
			*type_flag = t
			break
		}
	}
	if !valid {
		fmt.Fprintf(os.Stderr, "Invalid type: %q. Allowed values: %s\n", *type_flag, strings.Join(allowed_types, ", "))
		flag.Usage()
		os.Exit(1)
	}

	Create(*output_dir)

	if *type_flag == "track" {
		err := HandleTrack(ctx, *url_flag, *output_dir)
		if err != nil {
			log.Println(err)
		}
	}
	Clean()
}
