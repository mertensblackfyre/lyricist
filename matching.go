package main

import (
	"math"
	"strings"

	"github.com/lrstanley/go-ytdlp"
)

var FORBIDDEN_WORDS = [20]string{
	"bassboosted",
	"remix",
	"mv",
	"MV",
	"music video",
	"musicvideo",
	"video",
	"offcial",
	"remastered",
	"remaster",
	"reverb",
	"bassboost",
	"live",
	"acoustic",
	"8daudio",
	"concert",
	"acapella",
	"slowed",
	"instrumental",
	"cover",
}

func MatchScore(t *TrackDeezer, info *ytdlp.ExtractedInfo) int {
	score := 0

	if len(strings.TrimSpace(t.Title)) == len(strings.TrimSpace(*info.Title)) {
		score += 10
	}

	if !strings.Contains(strings.ToLower(*info.Title), strings.ToLower(t.Title)) {
		return 0
	}

	if !CheckForbiddenWords(*info.Title) {
		score += 5
	}
	if strings.Contains(strings.ToLower(*info.Channel), strings.ToLower(t.Artist.Name)) {
		score += 3
	}

	if strings.Contains(strings.ToLower(*info.Channel), "official") {
		score += 2
	}

	if int(math.Abs(float64(t.Duration)-*info.Duration)) < 3 {
		score += 5
	}

	return score
}

func CheckForbiddenWords(track_name string) bool {
	for _, word := range FORBIDDEN_WORDS {
		if strings.Contains(track_name, word) {
			return true
		}
	}
	return false
}
