package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func Transcode(id string, track *TrackDeezer, output string) error {
	coverFile := "%tempdir%" + "/" + id + ".jpg"
	audiofile := "%tempdir%" + "/" + id + ".webm"
	outputPath := filepath.Join(output, track.Title+".mp3")
	cmd := exec.Command("ffmpeg",
		"-i", audiofile,
		"-i", coverFile,
		"-c:a", "libmp3lame",
		"-b:a", "320k",
		"-metadata", "title="+track.Title,
		"-metadata", "artist="+track.Artist.Name,
		"-metadata", "album="+track.Album.Title,
		"-map", "0:a",
		"-map", "1:v",
		"-disposition:v", "attached_pic",
		"-y",
		outputPath,
	)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %w", err)
	}
	return nil
}
