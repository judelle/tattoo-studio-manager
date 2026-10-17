package appointment

import "errors"

var (
	ErrAppointmentNotFound         = errors.New("appointment not found")
	ErrAppointmentAlreadyCancelled = errors.New("appointment already cancelled")
	ErrAppointmentCompleted        = errors.New("appointment completed")
)
