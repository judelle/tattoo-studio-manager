package appointment

import (
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
