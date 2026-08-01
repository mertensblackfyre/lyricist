package main

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func HandlePlaylist(ctx context.Context, url string, output string) error {

	var concurrency = 2

	tracks_url, err := ExtractTrackURLCSV("everyday.csv")

	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Infof("Found %d tracks in the playlist", len(tracks_url))

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs = 0

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

			//logger.Printf("[%d/%d] Starting: %s", idx+1, len(tracks_url), track_)
			title, err := HandleTrack(ctx, track_, output)
			mu.Lock()
			if err != nil {
				errs++
			} else {
				logger.Infof("%s downloaded", title)
				success++
			}
			mu.Unlock()
		}(tracks, i)
	}

	wg.Wait()

	logger.Printf("Processed %d tracks: %d succeeded, %d failed", len(tracks_url), success, errs)

	return nil
}

func HandleTrack(ctx context.Context, url string, output string) (string, error) {

	track, err := ScrapeTrackInfo(ctx, url)
	logger.Infof("Processing %s", track.Title)

	if err != nil {
		return "", err
	}

	logger.Infof("Scraped Track Info - %s ", track.Title)
	t, err := GetTrackMetaData(ctx, &track)

	if err != nil {
		return "", err
	}
	logger.Infof("Fetched track's metadata from Deezer - %s ", track.Title)

	query := BuildSearchQuery(&t)
	info, err := Search(ctx, query)

	yt_duration := *info.Duration
	deez_duration := float64(t.Duration)

	if math.Abs(yt_duration-deez_duration) > 5 {
		logger.Warn("duration difference is more than 5")
	}

	id, err := DownloadTrack(ctx, *info.WebpageURL)
	if err != nil {
		return "", err
	}
	logger.Infof("Downloading track from youtube - %s ", track.Title)

	err = DownloadCoverImage(id, t.Album.CoverBig)
	if err != nil {
		logger.Warnf("error downloading cover image for %s: %s", t.Title, err)
	} else {
		logger.Infof("Downloaded track's cover image - %s ", track.Title)
	}

	err = Transcode(id, &t, output)
	if err != nil {
		return "", err
	}
	logger.Infof("Transcoded track - %s ", track.Title)
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
