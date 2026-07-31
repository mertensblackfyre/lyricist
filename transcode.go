package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"slices"
)

func Transcode(id string, track *TrackDeezer, output string) error {
	cover_file := "%tempdir%" + "/" + id + ".jpg"
	audio_file := "%tempdir%" + "/" + id + ".webm"
	output_path := filepath.Join(output, track.Artist.Name+" - "+track.Title+".mp3")

	if !Check(cover_file) {
		cover_file = ""
	}

	if !Check(audio_file) {
		return fmt.Errorf("ffmpeg error: audio does not exist")
	}

	args := []string{
		"-i", audio_file,
		"-c:a", "libmp3lame",
		"-b:a", "320k",
		"-metadata", "title=" + track.Title,
		"-metadata", "artist=" + track.Artist.Name,
		"-metadata", "album=" + track.Album.Title,
		"-map", "0:a",
		"-y",
		output_path,
	}

	if cover_file != "" {
		args = slices.Insert(args, 1, "-i", cover_file)
		args = append(args[:len(args)-1], "-map", "1:v", "-disposition:v", "attached_pic", output_path)
	}

	cmd := exec.Command("ffmpeg", args...)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %w", err)
	}
	return nil
}
