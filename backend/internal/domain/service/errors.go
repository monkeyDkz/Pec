package service

import "errors"

// Domain-level errors shared across services and surfaced by transport handlers.
var (
	ErrStreamNotLive = errors.New("stream is not live")
	ErrForbidden     = errors.New("forbidden")
	ErrNotFound      = errors.New("not found")
)
