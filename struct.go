package main

type Track struct {
	ID          string
	Title       string
	Artists     []string
	Album       string
	CoverArtURL string
	DurationMs  int
}

type trackJSON struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}
