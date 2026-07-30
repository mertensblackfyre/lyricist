# lyricist

Download Spotify tracks, albums, and playlists as tagged audio files.

## Features

- Downloads from Spotify track/album/playlist URLs
- Tags files with title, artist, album, and cover art (via Deezer)
- Outputs MP3 (320kbps) or any FFmpeg-supported format

## Dependencies

- [Go](https://go.dev/) 1.20+
- [yt-dlp](https://github.com/yt-dlp/yt-dlp)

## Installation

```bash
git clone https://github.com/mertensblackfyre/lyricist.git
cd lyricist
go build
```

## Usage

```bash
./lyricist -url "https://open.spotify.com/track/..." -type track
```

### Flags

- `-url`   : Spotify URL (required)
- `-type`  : `track`, `album`, or `playlist` (required)
- `-output`: Output directory (default: `output`)
