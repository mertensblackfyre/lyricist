package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func Transcode(id string, track *TrackDeezer, output string) error {
	cover_file := "%tempdir%" + "/" + id + ".jpg"
	audio_file := "%tempdir%" + "/" + id + ".webm"
	output_path := filepath.Join(output, track.Artist.Name+" - "+track.Title+".mp3")
	cmd := exec.Command("ffmpeg",
		"-i", audio_file,
		"-i", cover_file,
		"-c:a", "libmp3lame",
		"-b:a", "320k",
		"-metadata", "title="+track.Title,
		"-metadata", "artist="+track.Artist.Name,
		"-metadata", "album="+track.Album.Title,
		"-map", "0:a",
		"-map", "1:v",
		"-disposition:v", "attached_pic",
		"-y",
		output_path,
	)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %w", err)
	}
	return nil
}
