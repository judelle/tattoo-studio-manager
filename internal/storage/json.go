package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func Save(path string, snapshot Snapshot) error {
	if err := Validate(snapshot); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	tempFile, err := os.CreateTemp(filepath.Dir(path), "tattoo-studio-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp: %w", err)
	}
	tempPath := tempFile.Name()

	defer os.Remove(tempPath)

	_, err = tempFile.Write(data)
	if err != nil {
		tempFile.Close()
		return fmt.Errorf("write temp file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}

	err = os.Rename(tempPath, path)

	if err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}

	return nil
}

func Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return NewSnapshot(), nil
		}
		return Snapshot{}, fmt.Errorf("read json: %w", err)
	}
	var snapshot Snapshot
	err = json.Unmarshal(data, &snapshot)
	if err != nil {
		return Snapshot{}, fmt.Errorf("unmarshal json: %w", err)
	}
	if err := Validate(snapshot); err != nil {
		return Snapshot{}, fmt.Errorf("validate: %w", err)
	}
	return snapshot, nil
}
