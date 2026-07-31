package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
)

func HandlePlaylist(ctx context.Context, url string) {
	ScarapePlaylistTracks(ctx, url)
}
func HandleTrack(ctx context.Context, url string, output string) error {

	track, err := ScrapeTrackInfo(ctx, url)
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("error scraping %w", err)
	}
	t, err := GetTrackMetaData(ctx, &track)
	if err != nil {
		return fmt.Errorf("error fetching metadata from deezer %w", err)
	}
	query := BuildSearchQuery(&t)
	info, err := Search(ctx, query)

	yt_duration := *info.Duration
	deez_duration := float64(t.Duration)

	if math.Abs(yt_duration-deez_duration) > 5 {
		return fmt.Errorf("duration difference is more than 5, skipping\n")
	}

	id, err := DownloadTrack(ctx, *info.WebpageURL)
	if err != nil {
		return fmt.Errorf("error downloading track from youtube %w", err)
	}
	err = DownloadCoverImage(id, t.Album.CoverBig)
	if err != nil {
		fmt.Errorf("error downloading cover image for %s: %w", t.Title, err)
	}

	err = Transcode(id, &t, output)
	if err != nil {
		return fmt.Errorf("error transcoding %w", err)
	}
	return nil
}

func DownloadCoverImage(id string, url string) error {

	file_name := id + ".jpg"

	dir := "%tempdir%"
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to create directory: %v\n", err)
	}
	file_path := filepath.Join(dir, file_name)

	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to request image: %w\n", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned bad status: %s\n", response.Status)
	}

	file, err := os.Create(file_path)
	if err != nil {
		return fmt.Errorf("Failed to create file: %v\n", err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("Failed to save image: %v\n", err)
	}

	return nil
}
