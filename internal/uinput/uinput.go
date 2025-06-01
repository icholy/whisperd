package uinput

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/icholy/whisperd/internal/inputcodes"
)

const (
	UI_DEV_CREATE  = 0x5501
	UI_DEV_DESTROY = 0x5502
	UI_DEV_SETUP   = 0x405c5503
	UI_SET_EVBIT   = 0x40045564
	UI_SET_KEYBIT  = 0x40045565
)

type InputID struct {
	Bustype uint16
	Vendor  uint16
	Product uint16
	Version uint16
}

type Setup struct {
	ID           InputID
	Name         [80]byte
	FFEffectsMax uint32
}

func Create(name string) (*os.File, error) {
	if len(name) >= 80 {
		return nil, fmt.Errorf("name is too long: %q", name)
	}
	f, err := os.OpenFile("/dev/uinput", os.O_WRONLY|unix.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}
	fd := int(f.Fd())
	if err := unix.IoctlSetInt(fd, UI_SET_EVBIT, inputcodes.EV_KEY); err != nil {
		f.Close()
		return nil, err
	}
	for _, key := range inputcodes.Keys {
		if err := unix.IoctlSetInt(fd, UI_SET_KEYBIT, int(key)); err != nil {
			f.Close()
			return nil, err
		}
	}
	setup := Setup{
		ID: InputID{
			Bustype: 0x03, // USB
			Vendor:  0x1234,
			Product: 0x5678,
		},
	}
	copy(setup.Name[:], []byte(name))
	setup.Name[len(name)] = 0
	if err := unix.IoctlSetPointerInt(fd, UI_DEV_SETUP, int(uintptr(unsafe.Pointer(&setup)))); err != nil {
		f.Close()
		return nil, err
	}
	if err := unix.IoctlSetInt(fd, UI_DEV_CREATE, 0); err != nil {
		f.Close()
		return nil, err
	}
	return f, nil
}

func Destroy(f *os.File) error {
	fd := int(f.Fd())
	return unix.IoctlSetInt(fd, UI_DEV_DESTROY, 0)
}

func Emit(f *os.File, batch []inputcodes.Event) error {
	var buf bytes.Buffer
	for _, e := range batch {
		buf.Reset()
		if err := binary.Write(&buf, binary.LittleEndian, e); err != nil {
			return err
		}
		if _, err := f.Write(buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

func EmitText(f *os.File, text string) error {
	for i, r := range text {
		if i > 0 {
			// the output get scrambled without this delay
			time.Sleep(10 * time.Millisecond)
		}
		ee, ok := inputcodes.RuneEvents[r]
		if !ok {
			continue
		}
		batch := []inputcodes.Event{}
		// key down
		for _, e := range ee {
			e.Value = 1
			batch = append(batch, e)
		}
		batch = append(batch, inputcodes.Event{
			Type: inputcodes.EV_SYN,
			Code: inputcodes.SYN_REPORT,
		})
		// key up
		for _, e := range ee {
			e.Value = 0
			batch = append(batch, e)
		}
		batch = append(batch, inputcodes.Event{
			Type: inputcodes.EV_SYN,
			Code: inputcodes.SYN_REPORT,
		})
		if err := Emit(f, batch); err != nil {
			return err
		}
	}
	return nil
}
