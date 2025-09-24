package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type LogConfig struct {
	ID   string `json:"id"`
	Path string `json:"path"`
	Type string `json:"type"`
}

var ErrConfigFileAccess = errors.New("configuration file not accessible")

var ErrConfigParsing = errors.New("invalid configuration JSON")

type FileError struct {
	Path string
	Err  error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("%s: %v", e.Path, e.Err)
}

func (e *FileError) Unwrap() error { return e.Err }

func (e *FileError) Is(target error) bool {
	return target == ErrConfigFileAccess
}

type ParseError struct {
	Path string
	Err  error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s: %v", e.Path, e.Err)
}

func (e *ParseError) Unwrap() error { return e.Err }

func (e *ParseError) Is(target error) bool {
	return target == ErrConfigParsing
}

func Load(path string) ([]LogConfig, error) {
	if path == "" {
		return nil, &FileError{Path: path, Err: fmt.Errorf("%w: chemin vide", ErrConfigFileAccess)}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, &FileError{Path: path, Err: err}
	}

	var entries []LogConfig
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, &ParseError{Path: path, Err: err}
	}

	return entries, nil
}
