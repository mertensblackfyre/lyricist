package main

import (
	"encoding/json"
)

type Track struct {
	ID          string
	Title       string
	Artists     []string
	Album       string
	CoverArtURL string
	DurationMs  int
}

type trackJSON struct {
	Name     string          `json:"name"`
	ByArtist json.RawMessage `json:"byArtist"`
	Duration string          `json:"duration"`
	Image    string          `json:"image"`
	InAlbum  albumJSON       `json:"inAlbum"`
}

type albumJSON struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type artistJSON struct {
	Name string `json:"name"`
}
