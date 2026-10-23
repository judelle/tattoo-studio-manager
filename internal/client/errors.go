package client

import "errors"

var (
	ErrNameRequired       = errors.New("client name is empty")
	ErrPhoneRequired      = errors.New("client phone is empty")
	ErrPhoneAlreadyExists = errors.New("client phone already exists")
	ErrClientNotFound     = errors.New("client not found")
)
