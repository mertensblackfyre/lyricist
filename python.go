package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

func SearchYTMusic(query string) (map[string]any, error) {
	pythonCode := `
import sys
import json
from ytmusicapi import YTMusic

ytmusic = YTMusic()
results = ytmusic.search(sys.argv[1], filter='songs')
if results:
    print(json.dumps(results[0]))
else:
    print(json.dumps({}))
`

	cmd := exec.Command("python3", "-c", pythonCode, query)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		logger.Errorf("ytmusic error: %s", stderr.String())
		return nil, err
	}

	var result map[string]any
	json.Unmarshal(output, &result)
	return result, nil
}
