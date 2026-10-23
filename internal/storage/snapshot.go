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
