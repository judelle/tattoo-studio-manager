package appointment

import (
	"errors"
	"testing"
	"time"
)

func TestFindByID(t *testing.T) {
	appointments := []Appointment{
		{ID: "appointment-1"},
		{ID: "appointment-2"},
	}

	found, ok := FindByID(appointments, "appointment-2")

	if !ok {
		t.Fatal("expected appointment to be found")
	}

	if found.ID != "appointment-2" {
		t.Errorf(
			"expected ID appointment-2, got %s",
			found.ID,
		)
	}
}

func TestFindByIDNotFound(t *testing.T) {
	appointments := []Appointment{
		{ID: "appointment-1"},
	}

	_, ok := FindByID(appointments, "unknown")

	if ok {
		t.Error("expected appointment not to be found")
	}
}

func TestFindIndexByID(t *testing.T) {
	appointments := []Appointment{
		{ID: "appointment-1"},
		{ID: "appointment-2"},
		{ID: "appointment-3"},
	}

	index, ok := FindIndexByID(
		appointments,
		"appointment-2",
	)

	if !ok {
		t.Fatal("expected appointment to be found")
	}

	if index != 1 {
		t.Errorf(
			"expected index 1, got %d",
			index,
		)
	}
}

func TestFindIndexByIDNotFound(t *testing.T) {
	appointments := []Appointment{
		{ID: "appointment-1"},
	}

	index, ok := FindIndexByID(
		appointments,
		"unknown",
	)

	if ok {
		t.Error("expected appointment not to be found")
	}

	if index != -1 {
		t.Errorf(
			"expected index -1, got %d",
			index,
		)
	}
}

func TestFindByDate(t *testing.T) {
	location := time.FixedZone("test", 3*60*60)

	targetDate := time.Date(
		2026,
		time.October,
		20,
		0,
		0,
		0,
		0,
		location,
	)

	appointments := []Appointment{
		{
			ID: "late",
			StartAt: time.Date(
				2026,
				time.October,
				20,
				18,
				0,
				0,
				0,
				location,
			),
		},
		{
			ID: "another-day",
			StartAt: time.Date(
				2026,
				time.October,
				21,
				12,
				0,
				0,
				0,
				location,
			),
		},
		{
			ID: "early",
			StartAt: time.Date(
				2026,
				time.October,
				20,
				10,
				0,
				0,
				0,
				location,
			),
		},
	}

	result := FindByDate(
		appointments,
		targetDate,
	)

	if len(result) != 2 {
		t.Fatalf(
			"expected 2 appointments, got %d",
			len(result),
		)
	}

	if result[0].ID != "early" {
		t.Errorf(
			"expected first appointment early, got %s",
			result[0].ID,
		)
	}

	if result[1].ID != "late" {
		t.Errorf(
			"expected second appointment late, got %s",
			result[1].ID,
		)
	}
}

func TestFindByDateReturnsEmptySlice(t *testing.T) {
	targetDate := time.Date(
		2026,
		time.October,
		20,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	appointments := []Appointment{
		{
			ID: "1",
			StartAt: time.Date(
				2026,
				time.October,
				21,
				12,
				0,
				0,
				0,
				time.UTC,
			),
		},
	}

	result := FindByDate(
		appointments,
		targetDate,
	)

	if len(result) != 0 {
		t.Errorf(
			"expected empty result, got %d appointments",
			len(result),
		)
	}
}

func TestIsTimeTaken(t *testing.T) {
	startAt := time.Date(
		2026,
		time.October,
		20,
		15,
		0,
		0,
		0,
		time.UTC,
	)

	appointments := []Appointment{
		{
			ID:      "appointment-1",
			StartAt: startAt,
			Status:  StatusPlanned,
		},
	}

	taken := IsTimeTaken(
		appointments,
		startAt,
		"",
	)

	if !taken {
		t.Error("expected time to be taken")
	}
}

func TestIsTimeTakenIgnoresCancelledAppointment(t *testing.T) {
	startAt := time.Date(
		2026,
		time.October,
		20,
		15,
		0,
		0,
		0,
		time.UTC,
	)

	appointments := []Appointment{
		{
			ID:      "appointment-1",
			StartAt: startAt,
			Status:  StatusCancelled,
		},
	}

	taken := IsTimeTaken(
		appointments,
		startAt,
		"",
	)

	if taken {
		t.Error(
			"expected cancelled appointment not to occupy time",
		)
	}
}

func TestIsTimeTakenIgnoresCompletedAppointment(t *testing.T) {
	startAt := time.Date(
		2026,
		time.October,
		20,
		15,
		0,
		0,
		0,
		time.UTC,
	)

	appointments := []Appointment{
		{
			ID:      "appointment-1",
			StartAt: startAt,
			Status:  StatusCompleted,
		},
	}

	taken := IsTimeTaken(
		appointments,
		startAt,
		"",
	)

	if taken {
		t.Error(
			"expected completed appointment not to occupy time",
		)
	}
}

func TestIsTimeTakenIgnoresAppointmentByID(t *testing.T) {
	startAt := time.Date(
		2026,
		time.October,
		20,
		15,
		0,
		0,
		0,
		time.UTC,
	)

	appointments := []Appointment{
		{
			ID:      "appointment-1",
			StartAt: startAt,
			Status:  StatusPlanned,
		},
	}

	taken := IsTimeTaken(
		appointments,
		startAt,
		"appointment-1",
	)

	if taken {
		t.Error(
			"expected ignored appointment not to occupy time",
		)
	}
}

func TestHasAppointmentsForClient(t *testing.T) {
	appointments := []Appointment{
		{
			ID:       "appointment-1",
			ClientID: "client-1",
		},
		{
			ID:       "appointment-2",
			ClientID: "client-2",
		},
	}

	result := HasAppointmentsForClient(
		appointments,
		"client-2",
	)

	if !result {
		t.Error(
			"expected client to have appointments",
		)
	}
}

func TestHasAppointmentsForClientReturnsFalse(t *testing.T) {
	appointments := []Appointment{
		{
			ID:       "appointment-1",
			ClientID: "client-1",
		},
	}

	result := HasAppointmentsForClient(
		appointments,
		"client-999",
	)

	if result {
		t.Error(
			"expected client not to have appointments",
		)
	}
}

func TestFindByClientID(t *testing.T) {
	appointments := []Appointment{
		{
			ID:       "old",
			ClientID: "client-1",
			StartAt: time.Date(
				2026,
				time.January,
				10,
				12,
				0,
				0,
				0,
				time.UTC,
			),
		},
		{
			ID:       "other-client",
			ClientID: "client-2",
			StartAt: time.Date(
				2026,
				time.December,
				10,
				12,
				0,
				0,
				0,
				time.UTC,
			),
		},
		{
			ID:       "new",
			ClientID: "client-1",
			StartAt: time.Date(
				2026,
				time.October,
				10,
				12,
				0,
				0,
				0,
				time.UTC,
			),
		},
	}

	result := FindByClientID(
		appointments,
		"client-1",
	)

	if len(result) != 2 {
		t.Fatalf(
			"expected 2 appointments, got %d",
			len(result),
		)
	}

	if result[0].ID != "new" {
		t.Errorf(
			"expected newest appointment first, got %s",
			result[0].ID,
		)
	}

	if result[1].ID != "old" {
		t.Errorf(
			"expected oldest appointment second, got %s",
			result[1].ID,
		)
	}
}

func TestFindByClientIDReturnsEmptySlice(t *testing.T) {
	appointments := []Appointment{
		{
			ID:       "appointment-1",
			ClientID: "client-1",
		},
	}

	result := FindByClientID(
		appointments,
		"client-999",
	)

	if len(result) != 0 {
		t.Errorf(
			"expected empty result, got %d appointments",
			len(result),
		)
	}
}

func TestFindUpcoming(t *testing.T) {
	from := time.Date(
		2026,
		time.October,
		20,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	appointments := []Appointment{
		{
			ID: "future-late",
			StartAt: time.Date(
				2026,
				time.October,
				25,
				12,
				0,
				0,
				0,
				time.UTC,
			),
			Status: StatusPlanned,
		},
		{
			ID: "past",
			StartAt: time.Date(
				2026,
				time.October,
				19,
				12,
				0,
				0,
				0,
				time.UTC,
			),
			Status: StatusPlanned,
		},
		{
			ID: "future-near",
			StartAt: time.Date(
				2026,
				time.October,
				21,
				12,
				0,
				0,
				0,
				time.UTC,
			),
			Status: StatusPlanned,
		},
		{
			ID: "cancelled",
			StartAt: time.Date(
				2026,
				time.October,
				22,
				12,
				0,
				0,
				0,
				time.UTC,
			),
			Status: StatusCancelled,
		},
		{
			ID: "completed",
			StartAt: time.Date(
				2026,
				time.October,
				23,
				12,
				0,
				0,
				0,
				time.UTC,
			),
			Status: StatusCompleted,
		},
	}

	result := FindUpcoming(
		appointments,
		from,
	)

	if len(result) != 2 {
		t.Fatalf(
			"expected 2 upcoming appointments, got %d",
			len(result),
		)
	}

	if result[0].ID != "future-near" {
		t.Errorf(
			"expected nearest appointment first, got %s",
			result[0].ID,
		)
	}

	if result[1].ID != "future-late" {
		t.Errorf(
			"expected later appointment second, got %s",
			result[1].ID,
		)
	}
}

func TestFindUpcomingReturnsEmptySlice(t *testing.T) {
	from := time.Date(
		2026,
		time.October,
		20,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	appointments := []Appointment{
		{
			ID: "past",
			StartAt: time.Date(
				2026,
				time.October,
				19,
				12,
				0,
				0,
				0,
				time.UTC,
			),
			Status: StatusPlanned,
		},
		{
			ID: "cancelled",
			StartAt: time.Date(
				2026,
				time.October,
				21,
				12,
				0,
				0,
				0,
				time.UTC,
			),
			Status: StatusCancelled,
		},
	}

	result := FindUpcoming(
		appointments,
		from,
	)

	if len(result) != 0 {
		t.Errorf(
			"expected no upcoming appointments, got %d",
			len(result),
		)
	}
}

func TestCancel(t *testing.T) {
	tests := []struct {
		name           string
		appointments   []Appointment
		id             string
		expectedError  error
		expectedStatus Status
	}{
		{
			name: "planned appointment can be cancelled",
			appointments: []Appointment{
				{
					ID:     "1",
					Status: StatusPlanned,
				},
			},
			id:             "1",
			expectedError:  nil,
			expectedStatus: StatusCancelled,
		},
		{
			name: "cancelled appointment cannot be cancelled again",
			appointments: []Appointment{
				{
					ID:     "1",
					Status: StatusCancelled,
				},
			},
			id:             "1",
			expectedError:  ErrAppointmentCancelled,
			expectedStatus: StatusCancelled,
		},
		{
			name: "completed appointment cannot be cancelled",
			appointments: []Appointment{
				{
					ID:     "1",
					Status: StatusCompleted,
				},
			},
			id:             "1",
			expectedError:  ErrAppointmentCompleted,
			expectedStatus: StatusCompleted,
		},
		{
			name:          "appointment does not exist",
			appointments:  nil,
			id:            "999",
			expectedError: ErrAppointmentNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Cancel(
				tt.appointments,
				tt.id,
			)

			if !errors.Is(err, tt.expectedError) {
				t.Fatalf(
					"expected error %v, got %v",
					tt.expectedError,
					err,
				)
			}

			if len(tt.appointments) > 0 &&
				tt.appointments[0].Status != tt.expectedStatus {

				t.Errorf(
					"expected status %s, got %s",
					tt.expectedStatus,
					tt.appointments[0].Status,
				)
			}
		})
	}
}

func TestComplete(t *testing.T) {
	tests := []struct {
		name               string
		appointments       []Appointment
		id                 string
		finalPrice         int
		expectedError      error
		expectedFinalPrice int
		expectedStatus     Status
	}{
		{
			name: "planned appointment can be completed",
			appointments: []Appointment{
				{
					ID:         "1",
					Status:     StatusPlanned,
					FinalPrice: 0,
				},
			},
			id:                 "1",
			finalPrice:         5000,
			expectedError:      nil,
			expectedFinalPrice: 5000,
			expectedStatus:     StatusCompleted,
		},
		{
			name: "appointment can be completed with zero final price",
			appointments: []Appointment{
				{
					ID:         "1",
					Status:     StatusPlanned,
					FinalPrice: 0,
				},
			},
			id:                 "1",
			finalPrice:         0,
			expectedError:      nil,
			expectedFinalPrice: 0,
			expectedStatus:     StatusCompleted,
		},
		{
			name: "appointment not found",
			appointments: []Appointment{
				{
					ID:         "1",
					Status:     StatusPlanned,
					FinalPrice: 0,
				},
			},
			id:                 "999",
			finalPrice:         5000,
			expectedError:      ErrAppointmentNotFound,
			expectedFinalPrice: 0,
			expectedStatus:     StatusPlanned,
		},
		{
			name: "completed appointment cannot be completed again",
			appointments: []Appointment{
				{
					ID:         "1",
					Status:     StatusCompleted,
					FinalPrice: 5000,
				},
			},
			id:                 "1",
			finalPrice:         10000,
			expectedError:      ErrAppointmentCompleted,
			expectedFinalPrice: 5000,
			expectedStatus:     StatusCompleted,
		},
		{
			name: "cancelled appointment cannot be completed",
			appointments: []Appointment{
				{
					ID:         "1",
					Status:     StatusCancelled,
					FinalPrice: 0,
				},
			},
			id:                 "1",
			finalPrice:         5000,
			expectedError:      ErrAppointmentCancelled,
			expectedFinalPrice: 0,
			expectedStatus:     StatusCancelled,
		},
		{
			name: "negative final price is rejected",
			appointments: []Appointment{
				{
					ID:         "1",
					Status:     StatusPlanned,
					FinalPrice: 0,
				},
			},
			id:                 "1",
			finalPrice:         -5000,
			expectedError:      ErrFinalPriceLessThanZero,
			expectedFinalPrice: 0,
			expectedStatus:     StatusPlanned,
		},
		{
			name:          "empty appointments returns not found",
			appointments:  nil,
			id:            "999",
			finalPrice:    5000,
			expectedError: ErrAppointmentNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			completedAppointment, err := Complete(
				tt.appointments,
				tt.id,
				tt.finalPrice,
			)

			if !errors.Is(err, tt.expectedError) {
				t.Fatalf(
					"expected error %v, got %v",
					tt.expectedError,
					err,
				)
			}

			if err != nil {
				if len(tt.appointments) > 0 {
					if tt.appointments[0].FinalPrice != tt.expectedFinalPrice {
						t.Errorf(
							"expected final price %d, got %d",
							tt.expectedFinalPrice,
							tt.appointments[0].FinalPrice,
						)
					}

					if tt.appointments[0].Status != tt.expectedStatus {
						t.Errorf(
							"expected status %s, got %s",
							tt.expectedStatus,
							tt.appointments[0].Status,
						)
					}
				}

				return
			}

			if completedAppointment.FinalPrice != tt.expectedFinalPrice {
				t.Errorf(
					"expected final price %d, got %d",
					tt.expectedFinalPrice,
					completedAppointment.FinalPrice,
				)
			}

			if completedAppointment.Status != tt.expectedStatus {
				t.Errorf(
					"expected status %s, got %s",
					tt.expectedStatus,
					completedAppointment.Status,
				)
			}

			if len(tt.appointments) == 0 {
				t.Fatal("expected appointments slice not to be empty")
			}

			if tt.appointments[0].FinalPrice != tt.expectedFinalPrice {
				t.Errorf(
					"expected source appointment final price %d, got %d",
					tt.expectedFinalPrice,
					tt.appointments[0].FinalPrice,
				)
			}

			if tt.appointments[0].Status != tt.expectedStatus {
				t.Errorf(
					"expected source appointment status %s, got %s",
					tt.expectedStatus,
					tt.appointments[0].Status,
				)
			}
		})
	}
}

func TestReschedule(t *testing.T) {
	now := time.Date(
		2026,
		time.October,
		23,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	oldTime := now.Add(2 * time.Hour)
	pastTime := now.Add(-1 * time.Hour)
	futureTime := now.Add(24 * time.Hour)

	tests := []struct {
		name            string
		appointments    []Appointment
		id              string
		startAt         time.Time
		expectedError   error
		expectedStartAt time.Time
		expectedStatus  Status
	}{
		{
			name: "planned appointment can be rescheduled",
			appointments: []Appointment{
				{
					ID:      "1",
					Status:  StatusPlanned,
					StartAt: oldTime,
				},
			},
			id:              "1",
			startAt:         futureTime,
			expectedError:   nil,
			expectedStartAt: futureTime,
			expectedStatus:  StatusPlanned,
		},
		{
			name: "cancelled appointment becomes planned after reschedule",
			appointments: []Appointment{
				{
					ID:      "1",
					Status:  StatusCancelled,
					StartAt: oldTime,
				},
			},
			id:              "1",
			startAt:         futureTime,
			expectedError:   nil,
			expectedStartAt: futureTime,
			expectedStatus:  StatusPlanned,
		},
		{
			name: "completed appointment cannot be rescheduled",
			appointments: []Appointment{
				{
					ID:      "1",
					Status:  StatusCompleted,
					StartAt: oldTime,
				},
			},
			id:              "1",
			startAt:         futureTime,
			expectedError:   ErrAppointmentCompleted,
			expectedStartAt: oldTime,
			expectedStatus:  StatusCompleted,
		},
		{
			name: "appointment not found",
			appointments: []Appointment{
				{
					ID:      "1",
					Status:  StatusCancelled,
					StartAt: oldTime,
				},
			},
			id:              "999",
			startAt:         futureTime,
			expectedError:   ErrAppointmentNotFound,
			expectedStartAt: oldTime,
			expectedStatus:  StatusCancelled,
		},
		{
			name: "appointment cannot be rescheduled to the past",
			appointments: []Appointment{
				{
					ID:      "1",
					Status:  StatusCancelled,
					StartAt: oldTime,
				},
			},
			id:              "1",
			startAt:         pastTime,
			expectedError:   ErrStartTimeInPast,
			expectedStartAt: oldTime,
			expectedStatus:  StatusCancelled,
		},
		{
			name: "appointment cannot be rescheduled to occupied time",
			appointments: []Appointment{
				{
					ID:      "1",
					Status:  StatusCancelled,
					StartAt: oldTime,
				},
				{
					ID:      "2",
					Status:  StatusPlanned,
					StartAt: futureTime,
				},
			},
			id:              "1",
			startAt:         futureTime,
			expectedError:   ErrTimeTaken,
			expectedStartAt: oldTime,
			expectedStatus:  StatusCancelled,
		},
		{
			name:          "empty appointments returns not found",
			appointments:  []Appointment{},
			id:            "999",
			startAt:       futureTime,
			expectedError: ErrAppointmentNotFound,
		},
		{
			name: "appointment can be rescheduled to its own time",
			appointments: []Appointment{
				{
					ID:      "1",
					Status:  StatusPlanned,
					StartAt: futureTime,
				},
			},
			id:              "1",
			startAt:         futureTime,
			expectedError:   nil,
			expectedStartAt: futureTime,
			expectedStatus:  StatusPlanned,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rescheduleAt(
				tt.appointments,
				tt.id,
				tt.startAt,
				now,
			)

			if !errors.Is(err, tt.expectedError) {
				t.Fatalf(
					"expected error %v, got %v",
					tt.expectedError,
					err,
				)
			}

			if err != nil {
				if len(tt.appointments) > 0 {
					if !tt.appointments[0].StartAt.Equal(tt.expectedStartAt) {
						t.Errorf(
							"expected startAt %v, got %v",
							tt.expectedStartAt,
							tt.appointments[0].StartAt,
						)
					}

					if tt.appointments[0].Status != tt.expectedStatus {
						t.Errorf(
							"expected status %s, got %s",
							tt.expectedStatus,
							tt.appointments[0].Status,
						)
					}
				}

				return
			}

			if len(tt.appointments) == 0 {
				t.Fatal("expected appointments slice not to be empty")
			}

			if !tt.appointments[0].StartAt.Equal(tt.expectedStartAt) {
				t.Errorf(
					"expected source appointment startAt %v, got %v",
					tt.expectedStartAt,
					tt.appointments[0].StartAt,
				)
			}

			if tt.appointments[0].Status != tt.expectedStatus {
				t.Errorf(
					"expected source appointment status %s, got %s",
					tt.expectedStatus,
					tt.appointments[0].Status,
				)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	now := time.Date(
		2026,
		time.October,
		23,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	futureTime := now.Add(24 * time.Hour)
	pastTime := now.Add(-1 * time.Hour)

	tests := []struct {
		name           string
		appointments   []Appointment
		input          CreateInput
		expectedError  error
		expectedResult Appointment
	}{
		{
			name:         "valid appointment can be created",
			appointments: []Appointment{},
			input: CreateInput{
				ClientID:     "client-1",
				StartAt:      futureTime,
				Description:  "dragon",
				PlannedPrice: 10000,
				Deposit:      1500,
			},
			expectedError: nil,
			expectedResult: Appointment{
				ClientID:     "client-1",
				StartAt:      futureTime,
				Description:  "dragon",
				PlannedPrice: 10000,
				Deposit:      1500,
				Status:       StatusPlanned,
			},
		},
		{
			name:         "appointment cannot be created in the past",
			appointments: []Appointment{},
			input: CreateInput{
				ClientID:     "client-1",
				StartAt:      pastTime,
				Description:  "dragon",
				PlannedPrice: 10000,
				Deposit:      1500,
			},
			expectedError: ErrStartTimeInPast,
		},
		{
			name: "appointment cannot be created at occupied time",
			appointments: []Appointment{
				{
					ID:      "1",
					StartAt: futureTime,
					Status:  StatusPlanned,
				},
			},
			input: CreateInput{
				ClientID:     "client-1",
				StartAt:      futureTime,
				Description:  "dragon",
				PlannedPrice: 10000,
				Deposit:      1500,
			},
			expectedError: ErrTimeTaken,
		},
		{
			name: "cancelled appointment does not occupy time",
			appointments: []Appointment{
				{
					ID:      "1",
					StartAt: futureTime,
					Status:  StatusCancelled,
				},
			},
			input: CreateInput{
				ClientID:     "client-1",
				StartAt:      futureTime,
				Description:  "dragon",
				PlannedPrice: 10000,
				Deposit:      1500,
			},
			expectedError: nil,
			expectedResult: Appointment{
				ClientID:     "client-1",
				StartAt:      futureTime,
				Description:  "dragon",
				PlannedPrice: 10000,
				Deposit:      1500,
				Status:       StatusPlanned,
			},
		},
		{
			name: "completed appointment does not occupy time",
			appointments: []Appointment{
				{
					ID:      "1",
					StartAt: futureTime,
					Status:  StatusCompleted,
				},
			},
			input: CreateInput{
				ClientID:     "client-1",
				StartAt:      futureTime,
				Description:  "dragon",
				PlannedPrice: 10000,
				Deposit:      1500,
			},
			expectedError: nil,
			expectedResult: Appointment{
				ClientID:     "client-1",
				StartAt:      futureTime,
				Description:  "dragon",
				PlannedPrice: 10000,
				Deposit:      1500,
				Status:       StatusPlanned,
			},
		},
		{
			name:         "negative planned price is rejected",
			appointments: []Appointment{},
			input: CreateInput{
				ClientID:     "client-1",
				StartAt:      futureTime,
				Description:  "dragon",
				PlannedPrice: -10000,
				Deposit:      1500,
			},
			expectedError: ErrPlannedPriceLessThanZero,
		},
		{
			name:         "negative deposit is rejected",
			appointments: []Appointment{},
			input: CreateInput{
				ClientID:     "client-1",
				StartAt:      futureTime,
				Description:  "dragon",
				PlannedPrice: 10000,
				Deposit:      -1500,
			},
			expectedError: ErrDepositLessThanZero,
		},
		{
			name:         "zero planned price and zero deposit are allowed",
			appointments: []Appointment{},
			input: CreateInput{
				ClientID:     "client-1",
				StartAt:      futureTime,
				Description:  "dragon",
				PlannedPrice: 0,
				Deposit:      0,
			},
			expectedError: nil,
			expectedResult: Appointment{
				ClientID:     "client-1",
				StartAt:      futureTime,
				Description:  "dragon",
				PlannedPrice: 0,
				Deposit:      0,
				Status:       StatusPlanned,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdAppointment, err := createAt(
				tt.appointments,
				tt.input,
				now,
			)

			if !errors.Is(err, tt.expectedError) {
				t.Fatalf(
					"expected error %v, got %v",
					tt.expectedError,
					err,
				)
			}

			if err != nil {
				return
			}

			if createdAppointment.ID == "" {
				t.Error("expected generated ID, got empty string")
			}

			if createdAppointment.ClientID != tt.expectedResult.ClientID {
				t.Errorf(
					"expected client ID %s, got %s",
					tt.expectedResult.ClientID,
					createdAppointment.ClientID,
				)
			}

			if !createdAppointment.StartAt.Equal(tt.expectedResult.StartAt) {
				t.Errorf(
					"expected startAt %v, got %v",
					tt.expectedResult.StartAt,
					createdAppointment.StartAt,
				)
			}

			if createdAppointment.Description != tt.expectedResult.Description {
				t.Errorf(
					"expected description %s, got %s",
					tt.expectedResult.Description,
					createdAppointment.Description,
				)
			}

			if createdAppointment.PlannedPrice != tt.expectedResult.PlannedPrice {
				t.Errorf(
					"expected planned price %d, got %d",
					tt.expectedResult.PlannedPrice,
					createdAppointment.PlannedPrice,
				)
			}

			if createdAppointment.Deposit != tt.expectedResult.Deposit {
				t.Errorf(
					"expected deposit %d, got %d",
					tt.expectedResult.Deposit,
					createdAppointment.Deposit,
				)
			}

			if createdAppointment.Status != tt.expectedResult.Status {
				t.Errorf(
					"expected status %s, got %s",
					tt.expectedResult.Status,
					createdAppointment.Status,
				)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name            string
		appointments    []Appointment
		id              string
		wantErr         error
		wantLen         int
		wantRemainingID string
	}{
		{
			name: "delete planned appointment",
			appointments: []Appointment{
				{
					ID:     "1",
					Status: StatusPlanned,
				},
				{
					ID:     "2",
					Status: StatusPlanned,
				},
			},
			id:              "1",
			wantLen:         1,
			wantRemainingID: "2",
		},
		{
			name: "appointment not found",
			appointments: []Appointment{
				{
					ID:     "1",
					Status: StatusPlanned,
				},
			},
			id:      "999",
			wantErr: ErrAppointmentNotFound,
			wantLen: 1,
		},
		{
			name: "completed appointment cannot be deleted",
			appointments: []Appointment{
				{
					ID:     "1",
					Status: StatusCompleted,
				},
			},
			id:      "1",
			wantErr: ErrAppointmentCompleted,
			wantLen: 1,
		},
		{
			name: "cancelled appointment cannot be deleted",
			appointments: []Appointment{
				{
					ID:     "1",
					Status: StatusCancelled,
				},
			},
			id:      "1",
			wantErr: ErrAppointmentCancelled,
			wantLen: 1,
		},
		{
			name:         "empty appointments",
			appointments: []Appointment{},
			id:           "1",
			wantErr:      ErrAppointmentNotFound,
			wantLen:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Delete(
				tt.appointments,
				tt.id,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf(
					"Delete() error = %v, want %v",
					err,
					tt.wantErr,
				)
			}

			if len(got) != tt.wantLen {
				t.Fatalf(
					"Delete() len = %d, want %d",
					len(got),
					tt.wantLen,
				)
			}

			if tt.wantRemainingID != "" {
				if got[0].ID != tt.wantRemainingID {
					t.Errorf(
						"Delete() remaining ID = %q, want %q",
						got[0].ID,
						tt.wantRemainingID,
					)
				}
			}
		})
	}
}

func TestDeleteDoesNotModifyAppointmentOnError(t *testing.T) {
	appointments := []Appointment{
		{
			ID:         "1",
			Status:     StatusCompleted,
			FinalPrice: 5000,
		},
	}

	got, err := Delete(appointments, "1")

	if !errors.Is(err, ErrAppointmentCompleted) {
		t.Fatalf(
			"Delete() error = %v, want %v",
			err,
			ErrAppointmentCompleted,
		)
	}

	if len(got) != 1 {
		t.Fatalf(
			"Delete() len = %d, want 1",
			len(got),
		)
	}

	if got[0].ID != "1" {
		t.Errorf(
			"Delete() ID = %q, want %q",
			got[0].ID,
			"1",
		)
	}

	if got[0].Status != StatusCompleted {
		t.Errorf(
			"Delete() Status = %q, want %q",
			got[0].Status,
			StatusCompleted,
		)
	}

	if got[0].FinalPrice != 5000 {
		t.Errorf(
			"Delete() FinalPrice = %d, want 5000",
			got[0].FinalPrice,
		)
	}
}

func TestIsTimeTakenDoesNotIgnoreEmptyID(t *testing.T) {
	startAt := time.Date(
		2026,
		time.October,
		24,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	appointments := []Appointment{
		{
			ID:      "",
			StartAt: startAt,
			Status:  StatusPlanned,
		},
	}

	if !IsTimeTaken(appointments, startAt, "") {
		t.Fatal("expected time to be taken")
	}
}
