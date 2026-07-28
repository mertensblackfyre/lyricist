package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	ytdlp "github.com/lrstanley/go-ytdlp"
)

func Search(ctx context.Context, query string) ytdlp.ExtractedInfo {
	var builder strings.Builder

	builder.WriteString("ytsearch1:")
	builder.WriteString(query)

	result := builder.String()
	ytdlp.MustInstall(ctx, nil)

	dl := ytdlp.New()
	output, err := dl.Run(ctx, result, "--print-json", "--skip-download", "--no-playlist")
	if err != nil {
		panic(err)
	}
	var info ytdlp.ExtractedInfo
	if err := json.Unmarshal([]byte(output.Stdout), &info); err != nil {
		fmt.Errorf("error unmarshaling %w", err)
	}
	return info
}
