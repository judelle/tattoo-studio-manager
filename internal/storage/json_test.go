package storage

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"tattoo-studio/internal/appointment"
	"tattoo-studio/internal/client"
)

func TestSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "data.json")

	startAt := time.Date(
		2026,
		time.October,
		25,
		15,
		30,
		0,
		0,
		time.UTC,
	)

	snapshot := Snapshot{
		Version: 1,
		Clients: []client.Client{
			{
				ID:      "client-1",
				Name:    "Александр",
				IsAdult: true,
				Phone:   "+79990000000",
				Social:  "@alex",
			},
		},
		Appointments: []appointment.Appointment{
			{
				ID:           "appointment-1",
				ClientID:     "client-1",
				StartAt:      startAt,
				Description:  "dragon",
				PlannedPrice: 10000,
				Deposit:      1500,
				Status:       appointment.StatusPlanned,
			},
		},
	}

	err := Save(path, snapshot)
	if err != nil {
		t.Fatalf("Save() unexpected error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if loaded.Version != snapshot.Version {
		t.Errorf(
			"Version = %d, want %d",
			loaded.Version,
			snapshot.Version,
		)
	}

	if len(loaded.Clients) != 1 {
		t.Fatalf(
			"Clients len = %d, want 1",
			len(loaded.Clients),
		)
	}

	if loaded.Clients[0].ID != snapshot.Clients[0].ID {
		t.Errorf(
			"Client ID = %q, want %q",
			loaded.Clients[0].ID,
			snapshot.Clients[0].ID,
		)
	}

	if loaded.Clients[0].Name != snapshot.Clients[0].Name {
		t.Errorf(
			"Client Name = %q, want %q",
			loaded.Clients[0].Name,
			snapshot.Clients[0].Name,
		)
	}

	if loaded.Clients[0].Phone != snapshot.Clients[0].Phone {
		t.Errorf(
			"Client Phone = %q, want %q",
			loaded.Clients[0].Phone,
			snapshot.Clients[0].Phone,
		)
	}

	if len(loaded.Appointments) != 1 {
		t.Fatalf(
			"Appointments len = %d, want 1",
			len(loaded.Appointments),
		)
	}

	loadedAppointment := loaded.Appointments[0]
	expectedAppointment := snapshot.Appointments[0]

	if loadedAppointment.ID != expectedAppointment.ID {
		t.Errorf(
			"Appointment ID = %q, want %q",
			loadedAppointment.ID,
			expectedAppointment.ID,
		)
	}

	if loadedAppointment.ClientID != expectedAppointment.ClientID {
		t.Errorf(
			"Appointment ClientID = %q, want %q",
			loadedAppointment.ClientID,
			expectedAppointment.ClientID,
		)
	}

	if !loadedAppointment.StartAt.Equal(expectedAppointment.StartAt) {
		t.Errorf(
			"Appointment StartAt = %v, want %v",
			loadedAppointment.StartAt,
			expectedAppointment.StartAt,
		)
	}

	if loadedAppointment.Status != expectedAppointment.Status {
		t.Errorf(
			"Appointment Status = %q, want %q",
			loadedAppointment.Status,
			expectedAppointment.Status,
		)
	}
}

func TestSaveCreatesFile(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "data.json")

	snapshot := Snapshot{
		Version: 1,
	}

	err := Save(path, snapshot)
	if err != nil {
		t.Fatalf("Save() unexpected error: %v", err)
	}

	_, err = os.Stat(path)
	if err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestLoadNonExistingFile(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "does-not-exist.json")

	snapshot, err := Load(path)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if snapshot.Version != 0 {
		t.Errorf(
			"Version = %d, want 0",
			snapshot.Version,
		)
	}

	if snapshot.Clients != nil {
		t.Errorf(
			"Clients = %#v, want nil",
			snapshot.Clients,
		)
	}

	if snapshot.Appointments != nil {
		t.Errorf(
			"Appointments = %#v, want nil",
			snapshot.Appointments,
		)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "data.json")

	err := os.WriteFile(
		path,
		[]byte(`{"version":1,"clients":[`),
		0644,
	)
	if err != nil {
		t.Fatalf("failed to prepare test file: %v", err)
	}

	_, err = Load(path)
	if err == nil {
		t.Fatal("Load() expected error, got nil")
	}
}

func TestSaveOverwritesExistingData(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "data.json")

	first := Snapshot{
		Version: 1,
		Clients: []client.Client{
			{
				ID:    "1",
				Name:  "Старый клиент",
				Phone: "+79990000001",
			},
		},
	}

	err := Save(path, first)
	if err != nil {
		t.Fatalf("first Save() unexpected error: %v", err)
	}

	second := Snapshot{
		Version: 1,
		Clients: []client.Client{
			{
				ID:    "2",
				Name:  "Новый клиент",
				Phone: "+79990000002",
			},
		},
	}

	err = Save(path, second)
	if err != nil {
		t.Fatalf("second Save() unexpected error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if len(loaded.Clients) != 1 {
		t.Fatalf(
			"Clients len = %d, want 1",
			len(loaded.Clients),
		)
	}

	if loaded.Clients[0].ID != "2" {
		t.Errorf(
			"Client ID = %q, want %q",
			loaded.Clients[0].ID,
			"2",
		)
	}

	if loaded.Clients[0].Name != "Новый клиент" {
		t.Errorf(
			"Client Name = %q, want %q",
			loaded.Clients[0].Name,
			"Новый клиент",
		)
	}
}

func TestLoadReadError(t *testing.T) {
	tempDir := t.TempDir()

	_, err := Load(tempDir)
	if err == nil {
		t.Fatal("Load() expected error when path is directory, got nil")
	}

	if errors.Is(err, os.ErrNotExist) {
		t.Fatalf(
			"Load() returned ErrNotExist, expected another read error: %v",
			err,
		)
	}
}
