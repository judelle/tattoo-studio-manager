package appointment

import "errors"

var (
	ErrAppointmentNotFound    = errors.New("appointment not found")
	ErrAppointmentCompleted   = errors.New("appointment completed")
	ErrFinalPriceLessThanZero = errors.New("final price less than 0")
	ErrAppointmentCancelled   = errors.New("appointment cancelled")
)
