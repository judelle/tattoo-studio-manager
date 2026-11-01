package storage

import (
	"tattoo-studio/internal/appointment"
	"tattoo-studio/internal/client"
)

type Snapshot struct {
	Version      int                       `json:"version"`
	Clients      []client.Client           `json:"clients"`
	Appointments []appointment.Appointment `json:"appointments"`
}

const CurrentVersion = 1

func NewSnapshot() Snapshot {
	return Snapshot{
		Version:      CurrentVersion,
		Clients:      []client.Client{},
		Appointments: []appointment.Appointment{},
	}
}
