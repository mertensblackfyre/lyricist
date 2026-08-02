package main

import (
	"context"
	"fmt"
	"math"
	"sync"
)

func HandlePlaylist(ctx context.Context, file string, output string) error {

	var concurrency = 2

	tracks_url, err := ExtractTrackURLCSV(file)

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
			_, err := HandleTrack(ctx, track_, output)
			mu.Lock()
			if err != nil {
				errs++
			} else {
				//logger.Infof("%s downloaded", title)
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

	t = SanitizeDeezerTrack(&t)
	logger.Infof("Fetched & sanitized track's metadata from Deezer - %s ", track.Title)

	query := BuildSearchQuery(&t)
	info, err := Search(ctx, query)

	if err != nil {
		return "", err
	} else {
		yt_duration := *info.Duration
		deez_duration := float64(t.Duration)
		if math.Abs(yt_duration-deez_duration) > 5 {
			logger.Warn("duration difference is more than 5")
		}
	}

	logger.Infof("Downloading track from youtube - %s ", track.Title)
	id, err := DownloadTrack(ctx, *info.WebpageURL, output)
	if err != nil {
		return "", err
	}
	logger.Infof("Track downloaded - %s ", track.Title)
	RenameFileTrack(t.Artist.Name, t.Title, id, output)
	logger.Infof("Track file renamed - %s ", track.Title)

	/*
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
	*/
	return track.Title, nil
}
