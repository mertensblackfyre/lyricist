package main

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/lrstanley/go-ytdlp"
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
			_, err := HandleTrack(ctx, track_, output)
			mu.Lock()
			if err != nil {
				errs++
			} else {
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

	//t = SanitizeDeezerTrack(&t)
	logger.Infof("Fetched & sanitized track's metadata from Deezer - %s ", track.Title)

	if Check(filepath.Join(output, sanitize(track.Title)+".mp3")) {
		logger.Info("Track already exists, skipping")
		return track.Title, nil
	}

	t.Title = track.Title
	t.Artist.Name = track.Artists[0]
	query := BuildSearchQuery(&t)

	var info ytdlp.ExtractedInfo

	highest_score := 0

	var final ytdlp.ExtractedInfo
	for i := 1; i <= 5; i++ {

		info, err = Search(ctx, query, i)
		if err != nil {
			log.Error(err)
			break
		}

		score := MatchScore(&t, &info)
		log.Infof("Searched for %s --> score: %d", *info.Title, score)
		if score >= 20 {
			final = info
			break
		}
		if score > highest_score {
			final = info
			highest_score = score
		}
	}

	if final.ID == "" || *final.WebpageURL == "" {
		logger.Errorf("Track not found - %s", track.Title)
		return "", err
	}

	logger.Infof("Downloading track from youtube - %s ", track.Title)

	id, err := DownloadTrack(ctx, *final.WebpageURL, output)
	if err != nil {
		return "", err
	}
	logger.Infof("Track downloaded - %s ", track.Title)
	/*
		err = DownloadCoverImageDeezer(ctx, id, t.Album.CoverBig)
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
	RenameFileTrack(t.Artist.Name, t.Title, id, output)
	logger.Infof("Track file renamed - %s ", track.Title)
	return track.Title, nil
}
