package storage

import (
	"tattoo-studio/internal/appointment"
)

func Validate(snapshot Snapshot) error {
	if snapshot.Version != 1 {
		return ErrUnsupportedVersion
	}

	seenClientIDs := make(map[string]struct{})
	seenClientPhones := make(map[string]struct{})

	for _, currentClient := range snapshot.Clients {
		if currentClient.ID == "" {
			return ErrClientIDRequired
		}

		if _, exists := seenClientIDs[currentClient.ID]; exists {
			return ErrDuplicateClientID
		}

		seenClientIDs[currentClient.ID] = struct{}{}

		if currentClient.Name == "" {
			return ErrClientNameRequired
		}

		if currentClient.Phone == "" {
			return ErrClientPhoneRequired
		}

		if _, exists := seenClientPhones[currentClient.Phone]; exists {
			return ErrDuplicateClientPhone
		}

		seenClientPhones[currentClient.Phone] = struct{}{}
	}

	seenAppointmentIDs := make(map[string]struct{})
	plannedTimes := make(map[int64]struct{})

	for _, currentAppointment := range snapshot.Appointments {
		if currentAppointment.ID == "" {
			return ErrAppointmentIDRequired
		}

		if _, exists := seenAppointmentIDs[currentAppointment.ID]; exists {
			return ErrDuplicateAppointmentID
		}

		seenAppointmentIDs[currentAppointment.ID] = struct{}{}

		if _, exists := seenClientIDs[currentAppointment.ClientID]; !exists {
			return ErrAppointmentClientNotFound
		}

		switch currentAppointment.Status {
		case appointment.StatusPlanned,
			appointment.StatusCompleted,
			appointment.StatusCancelled:

		default:
			return ErrInvalidAppointmentStatus
		}

		if currentAppointment.StartAt.IsZero() {
			return ErrAppointmentStartAtRequired
		}

		if currentAppointment.PlannedPrice < 0 {
			return ErrAppointmentPlannedPriceLessThanZero
		}

		if currentAppointment.FinalPrice < 0 {
			return ErrAppointmentFinalPriceLessThanZero
		}

		if currentAppointment.Deposit < 0 {
			return ErrAppointmentDepositLessThanZero
		}

		if currentAppointment.Status == appointment.StatusPlanned {
			timeKey := currentAppointment.StartAt.UnixNano()

			if _, exists := plannedTimes[timeKey]; exists {
				return ErrDuplicatePlannedAppointmentTime
			}

			plannedTimes[timeKey] = struct{}{}
		}
	}

	return nil
}
