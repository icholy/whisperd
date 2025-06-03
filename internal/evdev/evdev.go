package evdev

import (
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
