package main

type TrackScrapeInfo struct {
	Title   string
	Artists []string
}

type TrackDeezer struct {
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	Md5Image string `json:"md5_image"`
	Artist   struct {
		Name       string `json:"name"`
		PictureBig string `json:"picture_big"`
	} `json:"artist"`
	Album struct {
		Title    string `json:"title"`
		CoverBig string `json:"cover_big"`
	} `json:"album"`
}
