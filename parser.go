package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func ParseRecordingJSON(htmlContent string) ([]string, error) {
	var result []string
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}
	var found_script string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" && n.FirstChild != nil {
			for _, attr := range n.Attr {
				if attr.Key == "type" && attr.Val == "application/ld+json" {
					if n.FirstChild != nil {
						found_script = n.FirstChild.Data
						return
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	if found_script == "" {
		return nil, fmt.Errorf("no JSON-LD script found")
	}

	r := make(map[string]any)
	json.Unmarshal([]byte(found_script), &r)

	if name, ok := r["name"].(string); ok {
		result = append(result, name)
	}
	if desc, ok := r["description"].(string); ok {
		result = append(result, desc)
	}
	if len(result) < 1 {
		return nil, fmt.Errorf("Something is missing")
	}
	return result, nil
}

func ParseArtists(description string) []string {
	parts := strings.Split(description, " · ")
	if len(parts) < 2 {
		return []string{"Unknown Artist"}
	}

	artist_part := parts[1]

	var artists []string
	for _, name := range strings.Split(artist_part, ",") {
		name = strings.TrimSpace(name)
		name = strings.Trim(name, `"`)
		if name != "" {
			artists = append(artists, name)
		}
	}

	if len(artists) == 0 {
		return []string{"Unknown Artist"}
	}
	return artists
}

func ParseSpotifyID(raw string) (string, error) {
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
