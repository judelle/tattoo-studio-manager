package storage

import "errors"

var (
	ErrUnsupportedVersion                  = errors.New("unsupported snapshot version")
	ErrClientIDRequired                    = errors.New("client ID is required")
	ErrDuplicateClientID                   = errors.New("duplicate client ID")
	ErrClientNameRequired                  = errors.New("client name is required")
	ErrClientPhoneRequired                 = errors.New("client phone is required")
	ErrDuplicateClientPhone                = errors.New("duplicate client phone")
	ErrAppointmentIDRequired               = errors.New("appointment ID is required")
	ErrDuplicateAppointmentID              = errors.New("duplicate appointment ID")
	ErrAppointmentClientNotFound           = errors.New("appointment client not found")
	ErrInvalidAppointmentStatus            = errors.New("invalid appointment status")
	ErrAppointmentPlannedPriceLessThanZero = errors.New("appointment planned price less than zero")
	ErrAppointmentFinalPriceLessThanZero   = errors.New("appointment final price less than zero")
	ErrAppointmentDepositLessThanZero      = errors.New("appointment deposit less than zero")
	ErrDuplicatePlannedAppointmentTime     = errors.New("duplicate planned appointment time")
	ErrAppointmentStartAtRequired          = errors.New("appointment start time is required")
)
