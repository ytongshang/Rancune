package util

import (
	"os"
	"path/filepath"
)

func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

func SelfDir() string {
	return filepath.Dir(SelfPath())
}

func MakeDirs(path string, perm os.FileMode) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, perm)
	}
	return nil
}
