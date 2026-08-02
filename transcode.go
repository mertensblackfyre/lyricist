package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
func Transcode(id string, track *TrackDeezer, output string) error {
	cover_file := filepath.Join("tmp", id+".jpg")
	audio_file := filepath.Join("tmp", id+".webm")
	output_path := filepath.Join(output, sanitize(track.Artist.Name+" - "+track.Title)+".mp3")

	info, err := os.Stat(audio_file)
	if err != nil || info == nil || info.Size() == 0 {
		logger.Errorf("audio file missing or empty: %s", audio_file)
		return err
	}

	args := []string{
		"-i", audio_file,
	}
	hasCover := false
	if _, err := os.Stat(cover_file); err == nil {
		args = append(args, "-i", cover_file)
		hasCover = true
	}

	args = append(args,
		"-c:a", "libmp3lame",
		"-b:a", "320k",
		"-id3v2_version", "3",
		"-metadata", fmt.Sprintf("title=%s", track.Title),
		"-metadata", fmt.Sprintf("artist=%s", track.Artist.Name),
		"-metadata", fmt.Sprintf("album=%s", track.Album.Title),
	)
	if hasCover {
		args = append(args,
			"-map", "0:a",
			"-map", "1:v",
			"-c:v", "copy",
			"-disposition:v", "attached_pic",
			"-metadata:s:v", "title=Album cover",
			"-metadata:s:v", "comment=Cover (front)",
		)
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
