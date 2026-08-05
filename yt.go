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
		"-f", "bestaudio/best",
		"-o", fmt.Sprintf("%s/%%(id)s.%%(ext)s", "tmp"),
		"--cookies-from-browser", "brave",
		"--user-agent", RandomUserAgent(),
		"--no-playlist",
		"--sleep-interval", "5",
		"--max-sleep-interval", "10",
		"--sleep-requests", "1",
		"--extractor-retries", "5",
		"--retries", "5",
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

func Search(ctx context.Context, query string, search int) (ytdlp.ExtractedInfo, error) {
	var builder strings.Builder

	builder.WriteString("ytsearch3:")
	builder.WriteString(query)

	result := builder.String()

	if err := limiter.Wait(ctx); err != nil {
		logger.Errorf("rate limit: %s", err)
		return ytdlp.ExtractedInfo{}, err
	}

	n := fmt.Sprintf("%d", search)

	output, err := dl.Run(ctx, result, "--print-json", "--skip-download", "--no-playlist", "--cookies-from-browser", "brave",
		"--user-agent", RandomUserAgent(),
		"--no-playlist",
		"--cookies-from-browser", "brave",
		"--sleep-interval", "5",
		"--max-sleep-interval", "10",
		"--sleep-requests", "1",
		"--extractor-retries", "5",
		"--retries", "5",
		"--playlist-items", n,
	)

	if err != nil {
		logger.Errorf("error running ytdlp %s", err)
		return ytdlp.ExtractedInfo{}, err
	}

	if len(output.Stdout) == 0 {
		logger.Errorf("yt-dlp returned no JSON (likely no results or rate-limited): %s", output.Stderr)
		return ytdlp.ExtractedInfo{}, fmt.Errorf(output.Stderr)
	}

	var info ytdlp.ExtractedInfo

	if err := json.Unmarshal([]byte(output.Stdout), &info); err != nil {
		logger.Errorf("error unmarshaling %s", err)
		return ytdlp.ExtractedInfo{}, err
	}

	return info, nil
}

func BuildSearchQuery(track *TrackDeezer) string {
	var builder strings.Builder
	builder.WriteString("\"")
	builder.WriteString(track.Artist.Name)
	builder.WriteString("\" \"")
	builder.WriteString(track.Title)
	builder.WriteString("\" official audio")
	url := builder.String()
	return url
}
