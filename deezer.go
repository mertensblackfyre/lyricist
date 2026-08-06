package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var base = "https://api.deezer.com/"

func GetTrackMetaData(ctx context.Context, track *TrackScrapeInfo) (TrackDeezer, error) {

	simple, complex := BuildQueryAPI(track)
	complex_body, err := SendRequest(ctx, complex, http.Client{
		Timeout: 15 * time.Second,
	})


	if err != nil {
		return TrackDeezer{}, err
	}

	var result DeezerSearchResponse
	if err := json.Unmarshal(complex_body, &result); err != nil {

		logger.Errorf("unmarshaling error of complex query: %s", err)
		return TrackDeezer{}, err
	}
	if len(result.Data) == 0 {
		logger.Warn("no data using complex query,trying simple query")

		simple_body, err := SendRequest(ctx, simple, http.Client{
			Timeout: 15 * time.Second,
		})

		if err != nil {
			return TrackDeezer{}, err
		}

		var result DeezerSearchResponse
		if err := json.Unmarshal(simple_body, &result); err != nil {
			logger.Errorf("unmarshaling error of simple query: %s", err)
			return TrackDeezer{}, err
		}
		if len(result.Data) == 0 {
			logger.Errorf("no data found using both queries")
			return TrackDeezer{}, err
		}

		t := result.Data[0]
		return t, nil
	}

	t := result.Data[0]
	return t, nil
}

func DownloadCoverImageDeezer(ctx context.Context, id string, url string) error {

	file_name := id + ".jpg"

	dir := "tmp"
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to create directory: %v\n", err)
	}
	file_path := filepath.Join(dir, file_name)

	body, err := SendRequest(ctx, url, http.Client{
		Timeout: 15 * time.Second,
	})

	file, err := os.Create(file_path)
	if err != nil {
		return fmt.Errorf("Failed to create file: %v\n", err)
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("Failed to save image: %v\n", err)
	}

	return nil
}

func BuildQueryAPI(track *TrackScrapeInfo) (string, string) {

	q := fmt.Sprintf(`"%s"`, CleanTrackTitle(track.Title))
	params := url.Values{}
	params.Set("q", q)
	simple_url := base + "search?" + params.Encode()

	q = fmt.Sprintf(`artist:"%s" track:"%s"`, track.Artists[0], CleanTrackTitle(track.Title))
	params = url.Values{}
	params.Set("q", q)

	complex_url := base + "search?" + params.Encode()

	return simple_url, complex_url

}
