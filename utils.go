package main

import (
	"errors"
	"fmt"
	"os"
)

func Create(output string) error {
	err := os.Mkdir(output, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("Failed to create single folder: %v\n", err)
	}
	return nil
}

func Check(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
func Clean() error {
	target_dir := "%tempdir%"
	err := os.RemoveAll(target_dir)
	if err != nil {
		return fmt.Errorf("Failed to delete folder: %v\n", err)
	}
	return nil
}
