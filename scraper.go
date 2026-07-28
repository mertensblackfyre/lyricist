package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Scraper struct {
	httpClient *http.Client
}

func NewScraper() *Scraper {
	return &Scraper{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (s *Scraper) GetTrack(ctx context.Context, url string) (Track, error) {

	id, err := extractTrackID(url)
	if err != nil {
		return Track{}, fmt.Errorf("parsing url error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Track{}, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; spotdl-bot/1.0)")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return Track{}, fmt.Errorf("fetching track page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Track{}, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Track{}, fmt.Errorf("reading response body: %w", err)
	}

	// Parse HTML and find JSON-LD script with MusicRecording type.
	trackJSON, err := extractMusicRecordingJSON(string(body))
	if err != nil {
		return Track{}, fmt.Errorf("extracting JSON-LD: %w", err)
	}

	durationMs, err := parseISO8601Duration(trackJSON.Duration)
	if err != nil {
		// Duration is not critical, we can default to 0 and continue.
		durationMs = 0
	}

	artists := parseArtists(trackJSON.ByArtist)

	// Prefer high-res cover art from album, fallback to track image.
	coverURL := trackJSON.Image
	if trackJSON.InAlbum.Image != "" {
		coverURL = trackJSON.InAlbum.Image
	}

	return Track{
		ID:          id,
		Title:       trackJSON.Name,
		Artists:     artists,
		Album:       trackJSON.InAlbum.Name,
		CoverArtURL: coverURL,
		DurationMs:  durationMs,
	}, nil
}

// extractMusicRecordingJSON parses the HTML and returns the MusicRecording JSON-LD data.
func extractMusicRecordingJSON(htmlContent string) (*trackJSON, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	var foundScript string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, attr := range n.Attr {
				if attr.Key == "type" && attr.Val == "application/ld+json" {
					if n.FirstChild != nil {
						foundScript = n.FirstChild.Data
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

	if foundScript == "" {
		return nil, fmt.Errorf("no JSON-LD script found")
	}

	// The JSON-LD may contain multiple items; try to unmarshal as array.
	var items []json.RawMessage
	if err := json.Unmarshal([]byte(foundScript), &items); err == nil {
		for _, raw := range items {
			var t trackJSON
			if err := json.Unmarshal(raw, &t); err == nil && t.Name != "" {
				return &t, nil
			}
		}
	}

	// Fallback: single item.
	var t trackJSON
	if err := json.Unmarshal([]byte(foundScript), &t); err != nil {
		return nil, fmt.Errorf("unmarshaling JSON-LD: %w", err)
	}
	if t.Name == "" {
		return nil, fmt.Errorf("track name missing in JSON-LD")
	}
	return &t, nil
}

// parseArtists handles the byArtist field which can be a single object or an array.
func parseArtists(raw json.RawMessage) []string {
	var single artistJSON
	if err := json.Unmarshal(raw, &single); err == nil {

		return []string{single.Name}
	}

	var multiple []artistJSON
	if err := json.Unmarshal(raw, &multiple); err == nil {
		names := make([]string, len(multiple))
		for i, a := range multiple {
			names[i] = a.Name
			fmt.Println(a.Name)
		}
		return names
	}

	return []string{"Unknown Artist"}
}

// parseISO8601Duration converts a duration string like "PT3M33S" to milliseconds.
func parseISO8601Duration(duration string) (int, error) {
	if duration == "" || !strings.HasPrefix(duration, "PT") {
		return 0, fmt.Errorf("invalid duration format: %s", duration)
	}

	d := duration[2:] // remove "PT"
	var h, m, s int
	var err error

	// Parse hours if present.
	if idx := strings.IndexByte(d, 'H'); idx != -1 {
		fmt.Sscanf(d[:idx], "%d", &h)
		d = d[idx+1:]
	}
	// Parse minutes if present.
	if idx := strings.IndexByte(d, 'M'); idx != -1 {
		fmt.Sscanf(d[:idx], "%d", &m)
		d = d[idx+1:]
	}
	// Parse seconds if present.
	if idx := strings.IndexByte(d, 'S'); idx != -1 {
		fmt.Sscanf(d[:idx], "%d", &s)
	}

	if h == 0 && m == 0 && s == 0 && err == nil {
		return 0, fmt.Errorf("no time components found in %q", duration)
	}

	return (h*3600 + m*60 + s) * 1000, nil
}
