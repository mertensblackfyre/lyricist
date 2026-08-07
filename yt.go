package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	ytdlp "github.com/lrstanley/go-ytdlp"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(rate.Every(3*time.Second), 2)
var dl = ytdlp.New()

func DownloadTrack(ctx context.Context, url string, output string) (string, error) {
	pos := strings.Index(url, "v=")

	if pos == -1 {
		logger.Errorf("can't extract id from youtube url")
		return "", nil
	}

	id := url[pos+2:]
	m_url := fmt.Sprintf("https://music.youtube.com/watch?v=%s", id)

	args := []string{
		m_url,
		"-f", "bestaudio[ext=m4a]",
		"-o", fmt.Sprintf("%s/%%(id)s.%%(ext)s", output),
		"--cookies-from-browser", "brave",
		"--user-agent", RandomUserAgent(),
		"--no-playlist",
		"--sleep-interval", "5",
		"--max-sleep-interval", "10",
		"--sleep-requests", "1",
		"--extractor-retries", "5",
		"--retries", "5",
		"--embed-metadata",
		"--embed-thumbnail",
	}

	if err := limiter.Wait(ctx); err != nil {
		logger.Errorf("rate limit: %s", err)
		return "", err
	}

	result, err := dl.Run(ctx, args...)
	if len(result.Stdout) == 0 {
		logger.Errorf("yt-dlp downloaded no song (likely no results or rate-limited): %s", result.Stderr)
		return "", fmt.Errorf(result.Stderr)
	}
	if err != nil {
		logger.Errorf("error running ytdlp %s", err)
		return "", err
	}
	return id, nil
}

func Search(ctx context.Context, query string, search int) ([]ytdlp.ExtractedInfo, error) {
	var builder strings.Builder

	builder.WriteString("ytsearch10:")
	builder.WriteString(query)

	result := builder.String()

	if err := limiter.Wait(ctx); err != nil {
		logger.Errorf("rate limit: %s", err)
		return []ytdlp.ExtractedInfo{}, err
	}

	output, err := dl.Run(ctx, result, "--print-json", "--skip-download", "--no-playlist", "--cookies-from-browser", "brave",
		"--user-agent", RandomUserAgent(),
		"--no-playlist",
		"--sleep-interval", "5",
		"--max-sleep-interval", "10",
		"--sleep-requests", "1",
		"--extractor-retries", "5",
		"--retries", "5",
		"--playlist-items", "1-10",
	)

	if err != nil {
		logger.Errorf("error running ytdlp %s", err)
		return []ytdlp.ExtractedInfo{}, err
	}

	if len(output.Stdout) == 0 {
		logger.Errorf("yt-dlp returned no JSON (likely no results or rate-limited): %s", output.Stderr)
		return []ytdlp.ExtractedInfo{}, fmt.Errorf(output.Stderr)
	}

	lines := strings.Split(strings.TrimSpace(output.Stdout), "\n")
	if len(lines) == 0 {
		logger.Errorf("no results found")
		return []ytdlp.ExtractedInfo{}, err
	}

	var results []ytdlp.ExtractedInfo
	for _, line := range lines {
		if line == "" {
			continue
		}
		var info ytdlp.ExtractedInfo
		err = json.Unmarshal([]byte(line), &info)
		if err != nil {
			logger.Warnf("error unmarshaling line: %s", err)
			continue
		}
		results = append(results, info)
	}

	if len(results) == 0 {
		logger.Errorf("no valid results parsed")
		return []ytdlp.ExtractedInfo{}, err
	}
	return results, nil
}

func BuildSearchQuery(track *TrackDeezer) string {
	var builder strings.Builder
	builder.WriteString(track.Artist.Name)
	builder.WriteString(" - ")
	builder.WriteString(track.Title)
	builder.WriteString(" audio")
	url := builder.String()
	return url
}
