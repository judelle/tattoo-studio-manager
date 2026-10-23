package studio

import (
	"errors"
	"testing"

	"tattoo-studio/internal/appointment"
	"tattoo-studio/internal/client"
)

func TestDeleteClient(t *testing.T) {
	tests := []struct {
		name         string
		clients      []client.Client
		appointments []appointment.Appointment
		clientID     string
		wantErr      error
		wantLen      int
	}{
		{
			name: "delete client without appointments",
			clients: []client.Client{
				{
					ID:   "1",
					Name: "Александр",
				},
				{
					ID:   "2",
					Name: "Анна",
				},
			},
			appointments: []appointment.Appointment{},
			clientID:     "1",
			wantLen:      1,
		},
		{
			name: "client not found",
			clients: []client.Client{
				{
					ID: "1",
				},
			},
			appointments: []appointment.Appointment{},
			clientID:     "999",
			wantErr:      client.ErrClientNotFound,
			wantLen:      1,
		},
		{
			name: "client has planned appointment",
			clients: []client.Client{
				{
					ID: "1",
				},
			},
			appointments: []appointment.Appointment{
				{
					ID:       "appointment-1",
					ClientID: "1",
					Status:   appointment.StatusPlanned,
				},
			},
			clientID: "1",
			wantErr:  ErrClientHasAppointments,
			wantLen:  1,
		},
		{
			name: "client has completed appointment",
			clients: []client.Client{
				{
					ID: "1",
				},
			},
			appointments: []appointment.Appointment{
				{
					ID:       "appointment-1",
					ClientID: "1",
					Status:   appointment.StatusCompleted,
				},
			},
			clientID: "1",
			wantErr:  ErrClientHasAppointments,
			wantLen:  1,
		},
		{
			name: "client has cancelled appointment",
			clients: []client.Client{
				{
					ID: "1",
				},
			},
			appointments: []appointment.Appointment{
				{
					ID:       "appointment-1",
					ClientID: "1",
					Status:   appointment.StatusCancelled,
				},
			},
			clientID: "1",
			wantErr:  ErrClientHasAppointments,
			wantLen:  1,
		},
		{
			name: "appointments belong to another client",
			clients: []client.Client{
				{
					ID: "1",
				},
				{
					ID: "2",
				},
			},
			appointments: []appointment.Appointment{
				{
					ID:       "appointment-1",
					ClientID: "2",
					Status:   appointment.StatusPlanned,
				},
			},
			clientID: "1",
			wantLen:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteClient(
				tt.clients,
				tt.appointments,
				tt.clientID,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"DeleteClient() error = %v, want %v",
					err,
					tt.wantErr,
				)
			}

			if len(got) != tt.wantLen {
				t.Fatalf(
					"DeleteClient() len = %d, want %d",
					len(got),
					tt.wantLen,
				)
			}

			if tt.wantErr != nil {
				return
			}

			for _, currentClient := range got {
				if currentClient.ID == tt.clientID {
					t.Errorf(
						"DeleteClient() client %q was not deleted",
						tt.clientID,
					)
				}
			}
		})
	}
}
