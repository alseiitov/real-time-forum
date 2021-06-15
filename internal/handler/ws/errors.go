package ws

import "errors"

var (
	errTooManyConnections = errors.New("too many connections")
	errInvalidEventType   = errors.New("invalid event type")
	errNoTokenReceived    = errors.New("no token received")
	errInvalidEventBody   = errors.New("invalid event body")
)
