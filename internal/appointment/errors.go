package appointment

import "errors"

var (
	ErrAppointmentNotFound      = errors.New("appointment not found")
	ErrAppointmentCompleted     = errors.New("appointment completed")
	ErrFinalPriceLessThanZero   = errors.New("final price less than 0")
	ErrAppointmentCancelled     = errors.New("appointment cancelled")
	ErrStartTimeInPast          = errors.New("start time is in the past")
	ErrTimeTaken                = errors.New("time is already taken")
	ErrPlannedPriceLessThanZero = errors.New("planned price less than 0")
	ErrDepositLessThanZero      = errors.New("deposit less than 0")
)
