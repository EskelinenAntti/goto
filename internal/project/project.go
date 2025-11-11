package project

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	Name       string
	WorkingDir string
}

var ErrNotADirectory = errors.New("Not a directory.")

func New(target string) (Project, error) {
	name := Name(target)

	fileInfo, err := os.Stat(target)
	if err != nil {
		return Project{}, err
	}
	if !fileInfo.IsDir() {
		return Project{}, fmt.Errorf("Invalid path %s: %w", target, ErrNotADirectory)
	}
	return Project{
		Name:       name,
		WorkingDir: target,
	}, nil
}

func Name(path string) string {
	basename := filepath.Base(path)
	sessionPrefix := strings.ReplaceAll(basename, ".", "_")
	return strings.Join([]string{sessionPrefix, hash(path)}, "-")
}

func hash(path string) string {
	hash := sha1.New()
	hash.Write([]byte(path))
	hashByteSlice := hash.Sum(nil)
	return fmt.Sprintf("%x", hashByteSlice)[:4]
}
