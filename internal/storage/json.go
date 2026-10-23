package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func Save(path string, snapshot Snapshot) error {
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}
	tempPath := path + ".tmp"

	err = os.WriteFile(tempPath, data, 0644)
	if err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}

	err = os.Rename(tempPath, path)
	if err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}

func Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Snapshot{}, nil
		}
		return Snapshot{}, fmt.Errorf("read json: %w", err)
	}
	var snapshot Snapshot
	err = json.Unmarshal(data, &snapshot)
	if err != nil {
		return Snapshot{}, fmt.Errorf("unmarshal json: %w", err)
	}
	return snapshot, nil
}
