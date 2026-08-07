package main

import (
	"math"
	"regexp"
	"strings"

	"github.com/lrstanley/go-ytdlp"
)

var forbidden_words = []string{
	"bassboosted", "bassboost", "remix", "mv", "music video", "musicvideo",
	"video", "offcial", "official video", "remastered", "remaster", "reverb",
	"live", "acoustic", "8daudio", "concert", "acapella", "slowed",
	"instrumental", "cover", "reaction", "parody", "opening", "full",
}

var clean_regex = regexp.MustCompile(`[^\w\s]`)

func MatchScore(dt *TrackDeezer, info *ytdlp.ExtractedInfo) int {
	if info == nil || info.Title == nil {
		return 0
	}

	yt_title := safe_string(info.Title)
	yt_channel := safe_string(info.Channel)
	yt_duration := safe_float(info.Duration)

	deezer_title_norm := normalize_string(dt.Title)
	yt_title_norm := normalize_string(yt_title)
	artist_norm := normalize_string(dt.Artist.Name)
	channel_norm := normalize_string(yt_channel)

	if !strings.Contains(yt_title_norm, deezer_title_norm) {
		return 0
	}

	score := 0

	if deezer_title_norm == yt_title_norm {
		score += 15
	} else if strings.HasPrefix(yt_title_norm, deezer_title_norm) {
		score += 8
	}

	if has_unwanted_forbidden_words(deezer_title_norm, yt_title_norm) {
		score -= 15
	} else {
		score += 5
	}

	if strings.Contains(channel_norm, artist_norm) {
		score += 5
	}

	if strings.HasSuffix(channel_norm, "topic") {
		score += 10
	} else if strings.Contains(channel_norm, "official") || strings.Contains(yt_title_norm, "official audio") {
		score += 4
	}

	if dt.Duration > 0 && yt_duration > 0 {
		diff := math.Abs(float64(dt.Duration) - yt_duration)

		switch {
		case diff <= 2:
			score += 10
		case diff <= 5:
			score += 5
		case diff <= 10:
			score += 2
		case diff > 20:
			score -= 10
		}
	}

	if score < 0 {
		return 0
	}

	return score
}

func has_unwanted_forbidden_words(deezer_title_norm, yt_title_norm string) bool {
	for _, word := range forbidden_words {
		if strings.Contains(yt_title_norm, word) {
			if !strings.Contains(deezer_title_norm, word) {
				return true
			}
		}
	}
	return false
}

func safe_string(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func safe_float(ptr *float64) float64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func normalize_string(s string) string {
	s = strings.ToLower(s)
	s = clean_regex.ReplaceAllString(s, " ")
	return strings.Join(strings.Fields(s), " ")
}
