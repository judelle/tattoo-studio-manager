package studio

import (
	"errors"
	"testing"
	"time"

	"tattoo-studio/internal/appointment"
	"tattoo-studio/internal/client"
)

func TestCreateAppointment(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name         string
		clients      []client.Client
		appointments []appointment.Appointment
		input        appointment.CreateInput
		wantErr      error
	}{
		{
			name: "existing client",
			clients: []client.Client{
				{
					ID:    "client-1",
					Name:  "Alex",
					Phone: "+79990000001",
				},
			},
			appointments: []appointment.Appointment{},
			input: appointment.CreateInput{
				ClientID:     "client-1",
				StartAt:      now.Add(time.Hour),
				Description:  "test",
				PlannedPrice: 5000,
				Deposit:      1000,
			},
			wantErr: nil,
		},
		{
			name: "client not found",
			clients: []client.Client{
				{
					ID:    "client-1",
					Name:  "Alex",
					Phone: "+79990000001",
				},
			},
			appointments: []appointment.Appointment{},
			input: appointment.CreateInput{
				ClientID:     "missing-client",
				StartAt:      now.Add(time.Hour),
				Description:  "test",
				PlannedPrice: 5000,
				Deposit:      1000,
			},
			wantErr: client.ErrClientNotFound,
		},
		{
			name: "appointment validation error",
			clients: []client.Client{
				{
					ID:    "client-1",
					Name:  "Alex",
					Phone: "+79990000001",
				},
			},
			appointments: []appointment.Appointment{},
			input: appointment.CreateInput{
				ClientID:     "client-1",
				StartAt:      now.Add(-time.Hour),
				Description:  "test",
				PlannedPrice: 5000,
				Deposit:      1000,
			},
			wantErr: appointment.ErrStartTimeInPast,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateAppointment(
				tt.clients,
				tt.appointments,
				tt.input,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"CreateAppointment() error = %v, want %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}
