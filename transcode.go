package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Transcode(id string, track *TrackDeezer, output string) error {
	cover_file := filepath.Join("tmp", id+".jpg")
	audio_file := filepath.Join("tmp", id+".webm")
	output_path := filepath.Join(output, track.Artist.Name+" - "+track.Title+".mp3")

	info, err := os.Stat(audio_file)
	if err != nil || info == nil || info.Size() == 0 {
		logger.Errorf("audio file missing or empty: %s", audio_file)
		return err
	}
	args := []string{
		"-i", audio_file,
		"-c:a", "libmp3lame",
		"-b:a", "320k",
		"-metadata", fmt.Sprintf("title=%s", track.Title),
		"-metadata", fmt.Sprintf("artist=%s", track.Artist.Name),
		"-metadata", fmt.Sprintf("album=%s", track.Album.Title),
	}

	if Check(cover_file) {
		args = append([]string{"-i", cover_file}, args...) // Add cover FIRST
		args = append(args, "-map", "1:v", "-disposition:v", "attached_pic")
	}

	args = append(args, "-y", "-loglevel", "debug", output_path)
	cmd := exec.Command("ffmpeg", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		logger.Errorf("ffmpeg error: %s", stderr.String())
		return err
	}
	return nil
}
