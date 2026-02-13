package tray

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

// Status represents the current state of whisperd.
type Status int

const (
	Idle Status = iota
	Recording
	Transcribing
)

var icons map[Status][]byte

func init() {
	icons = map[Status][]byte{
		Idle:         circleIcon(color.RGBA{128, 128, 128, 255}),
		Recording:    circleIcon(color.RGBA{220, 40, 40, 255}),
		Transcribing: circleIcon(color.RGBA{220, 200, 40, 255}),
	}
}

func circleIcon(c color.Color) []byte {
	const size = 22
	const center = size / 2
	const radius = 9
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := range size {
		for x := range size {
			dx := x - center
			dy := y - center
			if dx*dx+dy*dy <= radius*radius {
				img.Set(x, y, c)
			}
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

var tooltips = map[Status]string{
	Idle:         "whisperd: idle",
	Recording:    "whisperd: recording",
	Transcribing: "whisperd: transcribing",
}
