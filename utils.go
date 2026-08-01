package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
)

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
func Clean() error {
	target_dir := "%tempdir%"
	err := os.RemoveAll(target_dir)
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
