package appointment

import (
	"sort"
	"time"
)

func FindByID(appointments []Appointment, id string) (Appointment, bool) {
	for _, appointment := range appointments {
		if appointment.ID == id {
			return appointment, true
		}
	}

	return Appointment{}, false
}

func FindIndexByID(appointments []Appointment, id string) (int, bool) {
	for i, appointment := range appointments {
		if appointment.ID == id {
			return i, true
		}
	}

	return -1, false
}

func FindByDate(appointments []Appointment, startAt time.Time) []Appointment {
	result := make([]Appointment, 0)
	for _, appointment := range appointments {
		if appointment.StartAt.Day() == startAt.Day() &&
			appointment.StartAt.Month() == startAt.Month() &&
			appointment.StartAt.Year() == startAt.Year() {
			result = append(result, appointment)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].StartAt.Before(result[j].StartAt)
	})

	return result
}

func IsTimeTaken(appointments []Appointment, startAt time.Time) bool {
	for _, appointment := range appointments {
		if appointment.Status == StatusPlanned &&
			appointment.StartAt.Equal(startAt) {
			return true
		}
	}

	return false
}
