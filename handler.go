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
	"sync"
)

func HandlePlaylist(ctx context.Context, url string, output string) error {

	var concurrency = 2 // sensible default
	tracks_url, err := ScarapePlaylistTracks(ctx, url)

	log.Printf("Found %d", len(tracks_url))
	if err != nil {
		log.Println(err)
		return err
	}

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	success := 0

	for i, tracks := range tracks_url {
		select {
		case <-ctx.Done():
			wg.Wait()
			return fmt.Errorf("pipeline cancelled: %w", ctx.Err())
		default:
		}

		wg.Add(1)
		go func(track_ string, idx int) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }()

			log.Printf("[%d/%d] Starting: %s", idx+1, len(tracks_url), track_)

			title, err := HandleTrack(ctx, track_, output)
			mu.Lock()
			if err != nil {
				errs = append(errs, fmt.Errorf("track %s: %w", title, err))
			} else {
				success++
			}
			mu.Unlock()
		}(tracks, i)
	}

	wg.Wait()

	log.Printf("Processed %d tracks: %d succeeded, %d failed", len(tracks_url), success, len(errs))

	if len(errs) > 0 {
		return fmt.Errorf("some downloads failed")
	}

	return nil
}

func HandleTrack(ctx context.Context, url string, output string) (string, error) {

	track, err := ScrapeTrackInfo(ctx, url)
	log.Printf("Scraped Track Info - %s ", track.Title)
	if err != nil {
		log.Fatalln(err)
		return "", fmt.Errorf("error scraping %w", err)
	}
	t, err := GetTrackMetaData(ctx, &track)

	log.Printf("Getting Track's metadata from Deezer - %s ", track.Title)
	if err != nil {
		return "", fmt.Errorf("error fetching metadata from deezer %w", err)
	}
	query := BuildSearchQuery(&t)
	info, err := Search(ctx, query)

	yt_duration := *info.Duration
	deez_duration := float64(t.Duration)

	if math.Abs(yt_duration-deez_duration) > 5 {
		fmt.Errorf("duration difference is more than 5, skipping\n")
	}

	id, err := DownloadTrack(ctx, *info.WebpageURL)
	log.Printf("Downloading track from youtube - %s ", track.Title)
	if err != nil {
		return "", fmt.Errorf("error downloading track from youtube %w", err)
	}
	err = DownloadCoverImage(id, t.Album.CoverBig)

	log.Printf("Downloading track's cover image - %s ", track.Title)
	if err != nil {
		fmt.Errorf("error downloading cover image for %s: %w", t.Title, err)
	}

	log.Printf("Transcoding track - %s ", track.Title)
	err = Transcode(id, &t, output)
	if err != nil {
		return "", fmt.Errorf("error transcoding %w", err)
	}

	return track.Title, nil
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
