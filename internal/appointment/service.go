package appointment

import (
	"sort"
	"tattoo-studio/internal/datetime"
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

func IsTimeTaken(appointments []Appointment, startAt time.Time, ignoreID string) bool {
	for _, appointment := range appointments {
		if appointment.ID == ignoreID {
			continue
		}

		if appointment.Status == StatusPlanned &&
			appointment.StartAt.Equal(startAt) {
			return true
		}
	}

	return false
}

func HasAppointmentsForClient(appointments []Appointment, clientId string) bool {
	for _, appointment := range appointments {
		if appointment.ClientID == clientId {
			return true
		}
	}
	return false
}

func FindByClientID(appointments []Appointment, clientID string) []Appointment {
	result := make([]Appointment, 0)
	for _, v := range appointments {
		if v.ClientID == clientID {
			result = append(result, v)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].StartAt.After(result[j].StartAt)
	})
	return result
}

func FindUpcoming(appointments []Appointment, from time.Time) []Appointment {
	result := make([]Appointment, 0)
	for _, appointment := range appointments {
		if appointment.Status == StatusPlanned && from.Before(appointment.StartAt) {
			result = append(result, appointment)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].StartAt.Before(result[j].StartAt)
	})
	return result
}

func Cancel(appointments []Appointment, id string) error {
	index, found := FindIndexByID(appointments, id)
	if !found {
		return ErrAppointmentNotFound
	}

	switch appointments[index].Status {
	case StatusCancelled:
		return ErrAppointmentCancelled

	case StatusCompleted:
		return ErrAppointmentCompleted
	}

	appointments[index].Status = StatusCancelled

	return nil
}

func Complete(appointments []Appointment, id string, finalPrice int) (Appointment, error) {
	idx, found := FindIndexByID(
		appointments,
		id,
	)
	if !found {
		return Appointment{}, ErrAppointmentNotFound

	}

	if appointments[idx].Status == StatusCompleted {
		return Appointment{}, ErrAppointmentCompleted
	}

	if appointments[idx].Status == StatusCancelled {
		return Appointment{}, ErrAppointmentCancelled
	}

	if finalPrice < 0 {
		return Appointment{}, ErrFinalPriceLessThanZero
	}

	appointments[idx].FinalPrice = finalPrice
	appointments[idx].Status = StatusCompleted

	return appointments[idx], nil
}

func Reschedule(appointments []Appointment, id string, startAt time.Time) error {
	idx, ok := FindIndexByID(appointments, id)
	if !ok {
		return ErrAppointmentNotFound
	}

	if appointments[idx].Status == StatusCompleted {
		return ErrAppointmentCompleted
	}

	if startAt.Before(time.Now().In(datetime.Moscow)) {
		return ErrPastTimeReschedule
	}

	if IsTimeTaken(appointments, startAt, appointments[idx].ID) {
		return ErrTimeTaken
	}

	if appointments[idx].Status == StatusCancelled {
		appointments[idx].Status = StatusPlanned
	}

	appointments[idx].StartAt = startAt
	return nil
}
