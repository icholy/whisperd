package inputcodes

import "golang.org/x/sys/unix"

type Event struct {
	Time  unix.Timeval
	Type  uint16
	Code  uint16
	Value int32
}
