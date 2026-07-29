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

func (s *Scraper) ScrapeTrackInfo(ctx context.Context, url string) (Track, error) {

	id, err := ExtractTrackID(url)
	if err != nil {
		return Track{}, fmt.Errorf("parsing url error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Track{}, fmt.Errorf("creating request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Track{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Track{}, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Track{}, err
	}

	trackJSON, err := ExtractMusicRecordingJSON(string(body))
	if err != nil {
		return Track{}, fmt.Errorf("extracting JSON-LD: %w", err)
	}

	art := ParseArtistsFromDescription(trackJSON.Desc)

	return Track{
		ID:      id,
		Title:   trackJSON.Name,
		Artists: art,
	}, nil
}

func ExtractMusicRecordingJSON(htmlContent string) (*trackJSON, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}
	var foundScript string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" && n.FirstChild != nil {
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

	var items []json.RawMessage
	if err := json.Unmarshal([]byte(foundScript), &items); err == nil {
		for _, raw := range items {
			var t trackJSON
			if err := json.Unmarshal(raw, &t); err == nil && t.Name != "" {
				return &t, nil
			}
		}
	}

	var t trackJSON
	if err := json.Unmarshal([]byte(foundScript), &t); err != nil {
		return nil, fmt.Errorf("unmarshaling JSON-LD: %w", err)
	}
	if t.Name == "" {
		return nil, fmt.Errorf("track name missing in JSON-LD")
	}
	return &t, nil
}
