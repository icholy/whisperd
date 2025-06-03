package pipewire

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/icholy/whisperd/internal/wav"
)

// Options specifies the audio recording parameters for PipeWire.
type Options struct {
	SampleRate  int
	NumChannels int
}

// Recorder manages an audio recording session using PipeWire.
type Recorder struct {
	PCM     bytes.Buffer
	Process *os.Process
	Options Options
}

// Record starts a new audio recording session with the given options and returns a Recorder.
func Record(ctx context.Context, opt Options) (*Recorder, error) {
	var rec Recorder
	cmd := exec.CommandContext(ctx, "parec",
		"--format", "s16",
		"--rate", strconv.Itoa(opt.SampleRate),
		"--channels", strconv.Itoa(opt.NumChannels),
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = &rec.PCM
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	rec.Process = cmd.Process
	rec.Options = opt
	return &rec, nil
}

// Stop terminates the recording session and waits for the process to exit.
func (r *Recorder) Stop() error {
	if err := r.Process.Signal(syscall.SIGTERM); err != nil {
		return err
	}
	if _, err := r.Process.Wait(); err != nil {
		return err
	}
	return nil
}

// WriteWAV writes the recorded PCM data as a WAV file to the provided writer.
func (r *Recorder) WriteWAV(w io.Writer) error {
	return wav.Write(w, r.PCM.Bytes(), wav.Options{
		SampleRate:  r.Options.SampleRate,
		NumChannels: r.Options.NumChannels,
	})
}
