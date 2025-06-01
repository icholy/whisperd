package inputcodes

import "golang.org/x/sys/unix"

// Event represents a Linux input event, as defined in input-event-codes.h.
type Event struct {
	Time  unix.Timeval
	Type  uint16
	Code  uint16
	Value int32
}
