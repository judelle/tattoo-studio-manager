package storage

import (
	"errors"
	"testing"
	"time"

	"tattoo-studio/internal/appointment"
	"tattoo-studio/internal/client"
)

func validSnapshot() Snapshot {
	startAt := time.Date(
		2026,
		time.October,
		25,
		15,
		0,
		0,
		0,
		time.UTC,
	)

	return Snapshot{
		Version: 1,
		Clients: []client.Client{
			{
				ID:      "client-1",
				Name:    "Александр",
				IsAdult: true,
				Phone:   "+79990000001",
				Social:  "@alex",
			},
			{
				ID:      "client-2",
				Name:    "Анна",
				IsAdult: true,
				Phone:   "+79990000002",
				Social:  "@anna",
			},
		},
		Appointments: []appointment.Appointment{
			{
				ID:           "appointment-1",
				ClientID:     "client-1",
				StartAt:      startAt,
				Description:  "dragon",
				PlannedPrice: 10000,
				FinalPrice:   0,
				Deposit:      1500,
				Status:       appointment.StatusPlanned,
			},
		},
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		prepare func() Snapshot
		wantErr error
	}{
		{
			name: "valid snapshot",
			prepare: func() Snapshot {
				return validSnapshot()
			},
			wantErr: nil,
		},
		{
			name: "unsupported version",
			prepare: func() Snapshot {
				snapshot := validSnapshot()
				snapshot.Version = 2
				return snapshot
			},
			wantErr: ErrUnsupportedVersion,
		},
		{
			name: "client ID required",
			prepare: func() Snapshot {
				snapshot := validSnapshot()
				snapshot.Clients[0].ID = ""
				return snapshot
			},
			wantErr: ErrClientIDRequired,
		},
		{
			name: "duplicate client ID",
			prepare: func() Snapshot {
				snapshot := validSnapshot()
				snapshot.Clients[1].ID = snapshot.Clients[0].ID
				return snapshot
			},
			wantErr: ErrDuplicateClientID,
		},
		{
			name: "client name required",
			prepare: func() Snapshot {
				snapshot := validSnapshot()
				snapshot.Clients[0].Name = ""
				return snapshot
			},
			wantErr: ErrClientNameRequired,
		},
		{
			name: "client phone required",
			prepare: func() Snapshot {
				snapshot := validSnapshot()
				snapshot.Clients[0].Phone = ""
				return snapshot
			},
			wantErr: ErrClientPhoneRequired,
		},
		{
			name: "duplicate client phone",
			prepare: func() Snapshot {
				snapshot := validSnapshot()
				snapshot.Clients[1].Phone = snapshot.Clients[0].Phone
				return snapshot
			},
			wantErr: ErrDuplicateClientPhone,
		},
		{
			name: "appointment ID required",
			prepare: func() Snapshot {
				snapshot := validSnapshot()
				snapshot.Appointments[0].ID = ""
				return snapshot
			},
			wantErr: ErrAppointmentIDRequired,
		},
		{
			name: "duplicate appointment ID",
			prepare: func() Snapshot {
				snapshot := validSnapshot()

				secondAppointment := snapshot.Appointments[0]
				secondAppointment.ClientID = "client-2"
				secondAppointment.StartAt = secondAppointment.StartAt.Add(time.Hour)

				snapshot.Appointments = append(
					snapshot.Appointments,
					secondAppointment,
				)

				return snapshot
			},
			wantErr: ErrDuplicateAppointmentID,
		},
		{
			name: "appointment client not found",
			prepare: func() Snapshot {
				snapshot := validSnapshot()
				snapshot.Appointments[0].ClientID = "missing-client"
				return snapshot
			},
			wantErr: ErrAppointmentClientNotFound,
		},
		{
			name: "invalid appointment status",
			prepare: func() Snapshot {
				snapshot := validSnapshot()
				snapshot.Appointments[0].Status = appointment.Status("invalid")
				return snapshot
			},
			wantErr: ErrInvalidAppointmentStatus,
		},
		{
			name: "appointment start time required",
			prepare: func() Snapshot {
				snapshot := validSnapshot()
				snapshot.Appointments[0].StartAt = time.Time{}
				return snapshot
			},
			wantErr: ErrAppointmentStartAtRequired,
		},
		{
			name: "negative planned price",
			prepare: func() Snapshot {
				snapshot := validSnapshot()
				snapshot.Appointments[0].PlannedPrice = -1
				return snapshot
			},
			wantErr: ErrAppointmentPlannedPriceLessThanZero,
		},
		{
			name: "negative final price",
			prepare: func() Snapshot {
				snapshot := validSnapshot()
				snapshot.Appointments[0].FinalPrice = -1
				return snapshot
			},
			wantErr: ErrAppointmentFinalPriceLessThanZero,
		},
		{
			name: "negative deposit",
			prepare: func() Snapshot {
				snapshot := validSnapshot()
				snapshot.Appointments[0].Deposit = -1
				return snapshot
			},
			wantErr: ErrAppointmentDepositLessThanZero,
		},
		{
			name: "duplicate planned appointment time",
			prepare: func() Snapshot {
				snapshot := validSnapshot()

				snapshot.Appointments = append(
					snapshot.Appointments,
					appointment.Appointment{
						ID:           "appointment-2",
						ClientID:     "client-2",
						StartAt:      snapshot.Appointments[0].StartAt,
						Description:  "second",
						PlannedPrice: 5000,
						Deposit:      1000,
						Status:       appointment.StatusPlanned,
					},
				)

				return snapshot
			},
			wantErr: ErrDuplicatePlannedAppointmentTime,
		},
		{
			name: "cancelled appointment may use same time",
			prepare: func() Snapshot {
				snapshot := validSnapshot()

				snapshot.Appointments = append(
					snapshot.Appointments,
					appointment.Appointment{
						ID:           "appointment-2",
						ClientID:     "client-2",
						StartAt:      snapshot.Appointments[0].StartAt,
						PlannedPrice: 5000,
						Deposit:      1000,
						Status:       appointment.StatusCancelled,
					},
				)

				return snapshot
			},
			wantErr: nil,
		},
		{
			name: "completed appointment may use same time",
			prepare: func() Snapshot {
				snapshot := validSnapshot()

				snapshot.Appointments = append(
					snapshot.Appointments,
					appointment.Appointment{
						ID:           "appointment-2",
						ClientID:     "client-2",
						StartAt:      snapshot.Appointments[0].StartAt,
						PlannedPrice: 5000,
						FinalPrice:   5000,
						Deposit:      1000,
						Status:       appointment.StatusCompleted,
					},
				)

				return snapshot
			},
			wantErr: nil,
		},
		{
			name: "same instant in different timezones conflicts",
			prepare: func() Snapshot {
				snapshot := validSnapshot()

				firstTime := time.Date(
					2026,
					time.October,
					25,
					12,
					0,
					0,
					0,
					time.UTC,
				)

				moscow := time.FixedZone("MSK", 3*60*60)

				sameInstantInMoscow := time.Date(
					2026,
					time.October,
					25,
					15,
					0,
					0,
					0,
					moscow,
				)

				snapshot.Appointments[0].StartAt = firstTime

				snapshot.Appointments = append(
					snapshot.Appointments,
					appointment.Appointment{
						ID:           "appointment-2",
						ClientID:     "client-2",
						StartAt:      sameInstantInMoscow,
						PlannedPrice: 5000,
						Deposit:      1000,
						Status:       appointment.StatusPlanned,
					},
				)

				return snapshot
			},
			wantErr: ErrDuplicatePlannedAppointmentTime,
		},
		{
			name: "empty clients and appointments are valid",
			prepare: func() Snapshot {
				return Snapshot{
					Version:      1,
					Clients:      []client.Client{},
					Appointments: []appointment.Appointment{},
				}
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshot := tt.prepare()

			err := Validate(snapshot)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"Validate() error = %v, want %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}
