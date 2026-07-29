package main

import (
	"fmt"
	"net/url"
	"strings"
)

func ExtractTrackID(raw string) (string, error) {
	if strings.HasPrefix(raw, "https://") || strings.HasPrefix(raw, "http://") {
		u, err := url.Parse(raw)
		if err != nil {
			return "", fmt.Errorf("invalid URL: %w", err)
		}
		parts := strings.Split(strings.TrimPrefix(u.Path, "/"), "/")
		if len(parts) < 2 || parts[0] != "track" {
			return "", fmt.Errorf("URL doesn't contain track ID")
		}
		id := strings.Split(parts[1], "?")[0]
		return id, nil
	}
	if strings.HasPrefix(raw, "spotify:") {
		parts := strings.Split(raw, ":")
		if len(parts) >= 3 && parts[1] == "track" {
			return parts[2], nil
		}
		return "", fmt.Errorf("invalid Spotify URI")
	}
	if len(raw) == 22 {
		return raw, nil
	}
	return "", fmt.Errorf("unrecognized track identifier: %s", raw)
}

func ParseArtistsFromDescription(description string) []string {
	parts := strings.Split(description, " · ")
	if len(parts) < 2 {
		return []string{"Unknown Artist"}
	}

	// The artist(s) are in parts[1], e.g. "Eminem, Royce Da 5'9\""
	artistPart := parts[1]

	var artists []string
	for _, name := range strings.Split(artistPart, ",") {
		name = strings.TrimSpace(name)
		name = strings.Trim(name, `"`) // remove any stray quotes
		if name != "" {
			artists = append(artists, name)
		}
	}

	if len(artists) == 0 {
		return []string{"Unknown Artist"}
	}
	return artists
}

func BuildURL(track *Track) string {

	var builder strings.Builder
	var has_unknown bool = false

	for _, artist := range track.Artists {
		if artist == "Unknown Artist" {
			has_unknown = true
			break
		}
		builder.WriteString(artist)
		builder.WriteString(", ")
	}

	if !has_unknown {
		builder.WriteString(" - ")
	}

	builder.WriteString(track.Title)
	builder.WriteString("audio")
	url := builder.String()
	return url
}
