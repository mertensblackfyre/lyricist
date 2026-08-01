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
