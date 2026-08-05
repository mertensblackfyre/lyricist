package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	ctx := context.Background()
	start := time.Now()

	allowed_types := []string{"track", "album", "playlist"}

	url_flag := flag.String("url", "", "Spotify URL to process")
	file_flag := flag.String("file", "", "File path to process")
	type_flag := flag.String("type", "", fmt.Sprintf("Type of resource (%s)", strings.Join(allowed_types, ", ")))
	output_dir := flag.String("output", "output", "Output directory for downloaded files")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: lyricist -url <spotify_url> | -file <path> -type <type> [-output <dir>]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if (*url_flag == "" && *file_flag == "") || *type_flag == "" {
		flag.Usage()
		os.Exit(1)
	}

	lower_type := strings.ToLower(*type_flag)
	valid := false
	for _, t := range allowed_types {
		if lower_type == t {
			valid = true
			*type_flag = t
			break
		}
	}
	if !valid {
		fmt.Fprintf(os.Stderr, "Invalid type: %q\n", *type_flag)
		flag.Usage()
		os.Exit(1)
	}

	Clean("tmp")
	Create(*output_dir)
	Create("tmp")

	if *type_flag == "track" {
		HandleTrack(ctx, *url_flag, *output_dir)
	}

	if *type_flag == "playlist" {
		HandlePlaylist(ctx, *file_flag, *output_dir)
	}
	logger.Infof("\n\nTotal execution time: %v", time.Since(start))
}
