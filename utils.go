package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func sanitize(s string) string {
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "\\", "-")
	s = strings.ReplaceAll(s, ":", "-")
	s = strings.ReplaceAll(s, "?", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "|", "-")
	s = strings.ReplaceAll(s, "<", "")
	s = strings.ReplaceAll(s, ">", "")
	return s
}

func RenameFileTrack(artist string, title string, id string, output string) {

	current := filepath.Join(output, id+".m4a")
	sanitized := sanitize(title)
	final := filepath.Join(output, sanitized+".m4a")
	os.Rename(current, final)
}

func Create(output string) error {
	err := os.Mkdir(output, 0755)
	if err != nil && !os.IsExist(err) {
		logger.Errorf("Failed to create single folder: %v\n", err)
		return err
	}
	return nil
}

func Check(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
func Clean(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		logger.Errorf("Failed to delete folder: %v\n", err)
		return err
	}
	return nil
}

func ExtractTrackURLCSV(filePath string) ([]string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		logger.Errorf("opening CSV: %s", err)
		return nil, err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	var urls []string
	lineNum := 0

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		lineNum++
		if lineNum == 1 {
			continue
		}

		if len(record) == 0 {
			continue
		}

		uri := record[0]
		const prefix = "spotify:track:"
		if !strings.HasPrefix(uri, prefix) {
			fmt.Fprintf(os.Stderr, "skipping row %d: invalid URI %q\n", lineNum, uri)
			continue
		}
		id := strings.TrimPrefix(uri, prefix)
		url := "https://open.spotify.com/track/" + id
		urls = append(urls, url)
	}

	return urls, nil
}

func CleanTrackTitle(title string) string {
	reParen := regexp.MustCompile(`(?i)\s*\(.*?(feat|ft|featuring).*?\)`)
	title = reParen.ReplaceAllString(title, "")

	re_bracket := regexp.MustCompile(`(?i)\s*\[.*?(official (audio|video|music video)|lyrics|hd|clean|explicit).*?\]`)
	title = re_bracket.ReplaceAllString(title, "")

	re_suffix := regexp.MustCompile(`(?i)\s*[-–]\s*(single|remaster(ed)?|deluxe edition|album version).*$`)
	title = re_suffix.ReplaceAllString(title, "")

	return strings.TrimSpace(title)
}
