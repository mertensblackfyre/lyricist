package main

import (
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Transcode(track *TrackDeezer) {
	err := ffmpeg.Input(tempAudioPath).
		Input(coverImagePath).
		Output(outputPath, ffmpeg.KwArgs{
			"c:a": "libmp3lame", // MP3 codec (or "libopus", "aac")
			"b:a": "320k",       // bitrate
			"metadata": []string{
				"title=" + track.Title,
				"artist=" + strings.Join(track.Artist.Name, ", "),
				"album=" + track.Album.Title,
			},
			"map":           []string{"0:a", "1:v"}, // audio from 1st input, image from 2nd
			"disposition:v": "attached_pic",         // embed as cover (for MP3)
		}).
		OverWriteOutput().ErrorToStdOut().Run()
}
