package main

import (
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

func Clean() error {
	targetDir := "%tempdir%"
	err := os.RemoveAll(targetDir)
	if err != nil {
		return fmt.Errorf("Failed to delete folder: %v\n", err)
	}
	return nil
}
