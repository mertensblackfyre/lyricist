package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func Transcode(id string, track *TrackDeezer, output string) error {
	cover_file := filepath.Join("tmp", id+".jpg")
	audio_file := filepath.Join("tmp", id+".mp3")
	output_path := filepath.Join(output, track.Artist.Name+" - "+track.Title+".mp3")

	if !Check(cover_file) {
		cover_file = ""
	}

	info, err := os.Stat(audio_file)
	if err != nil || info == nil || info.Size() == 0 {
		logger.Errorf("audio file missing or empty: %s", audio_file)
		return err
	}
	args := []string{
		"-i", audio_file,
		"-c:a", "libmp3lame",
		"-b:a", "320k",
		"-metadata", "title=" + track.Title,
		"-metadata", "artist=" + track.Artist.Name,
		"-metadata", "album=" + track.Album.Title,
		"-map", "0:a",
	}

	if cover_file != "" {
		args = append([]string{"-i", cover_file}, args...) // Add cover FIRST
		args = append(args, "-map", "1:v", "-disposition:v", "attached_pic")
	}

	args = append(args, "-y", "-loglevel", "debug", output_path)
	cmd := exec.Command("ffmpeg", args...)

	err = cmd.Run()
	if err != nil {
		logger.Errorf("ffmpeg error: %s", err)
		return err
	}

	return nil
}
