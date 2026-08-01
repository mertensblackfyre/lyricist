lyricist
============
Download Spotify tracks as tagged audio files.The tool could also download playlists, but the playlist should exported to a csv. There are online tools for that. This project is tailored for Exportify specifically

Dependencies
-----------
* Go 1.20+
* yt-dlp

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
