lyricist
============
Download Spotify tracks as tagged audio files.The tool could also download playlists, but the playlist should be exported to a CSV. This project is specifically designed for use with Exportify.

Dependencies
-----------
* Go 1.20+
* yt-dlp
* ytmusicapi
* python 3+

Installation
-----------
```
git clone https://github.com/mertensblackfyre/lyricist.git
cd lyricist
make build
./lyricist
```

Usage
-----------

For individual tracks:
```
./lyricist -url "https://open.spotify.com/track/..." -type track
```

For exported playlists:
```
./lyricist -file playlist.csv -type playlist
```
