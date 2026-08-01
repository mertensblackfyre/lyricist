package main

type TrackScrapeInfo struct {
	Title   string
	Artists []string
}

type TrackDeezer struct {
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	Artist   struct {
		Name string `json:"name"`
	} `json:"artist"`
	Album struct {
		Title    string `json:"title"`
		CoverBig string `json:"cover_big"`
	} `json:"album"`
}
type DeezerSearchResponse struct {
    Data  []TrackDeezer `json:"data"`
    Total int           `json:"total"`
}
func SanitizeDeezerTrack(d *TrackDeezer) TrackDeezer {
	if d == nil {
		return TrackDeezer{
			Title:    "Unknown Title",
			Duration: 0,
			Artist: struct {
				Name string `json:"name"`
			}{Name: "Unknown Artist"},
			Album: struct {
				Title    string `json:"title"`
				CoverBig string `json:"cover_big"`
			}{Title: "Unknown Album"},
		}
	}

	safe := *d
	if safe.Title == "" {
		safe.Title = "Unknown Title"
	}
	if safe.Artist.Name == "" {
		safe.Artist.Name = "Unknown Artist"
	}
	if safe.Album.Title == "" {
		safe.Album.Title = "Unknown Album"
	}
	return safe
}
