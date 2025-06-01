package evdev

import (
	"context"
	"encoding/binary"
	"os"

	"github.com/icholy/whisperd/internal/inputcodes"
)

// WaitForKey blocks until the specified key event (code and value) is received from the device.
func WaitForKey(device *os.File, code uint16, value int32) error {
	for {
		var e inputcodes.Event
		if err := binary.Read(device, binary.LittleEndian, &e); err != nil {
			return err
		}
		if e.Type == inputcodes.EV_KEY && e.Code == code && e.Value == value {
			return nil
		}
	}
}

// KeyDownContext returns a context that is canceled when the specified key is released after being pressed.
func KeyDownContext(device *os.File, code uint16) (context.Context, error) {
	if err := WaitForKey(device, code, 1); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancelCause(context.Background())
	go func() {
		cancel(WaitForKey(device, code, 0))
	}()
	return ctx, nil
}
