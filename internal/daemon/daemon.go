package daemon

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/icholy/whisperd/internal/evdev"
	"github.com/icholy/whisperd/internal/openai"
	"github.com/icholy/whisperd/internal/pipewire"
	"github.com/icholy/whisperd/internal/tray"
	"github.com/icholy/whisperd/internal/uinput"
)

type Daemon struct {
	Input   *os.File
	Output  *os.File
	Client  openai.Client
	KeyCode uint16
	Dump    bool
}

func (d *Daemon) Run(ctx context.Context) error {
	for {
		tray.SetStatus(tray.Idle)
		log.Println("waiting for key down")
		if err := evdev.WaitForKey(d.Input, d.KeyCode, 1); err != nil {
			return fmt.Errorf("wait for key down: %w", err)
		}
		tray.SetStatus(tray.Recording)
		log.Println("starting recording")
		rec, err := pipewire.Record(ctx, pipewire.Options{
			SampleRate:  16000,
			NumChannels: 1,
		})
		if err != nil {
			return fmt.Errorf("start recording: %w", err)
		}
		log.Println("waiting for key up")
		if err := evdev.WaitForKey(d.Input, d.KeyCode, 0); err != nil {
			return fmt.Errorf("wait for key up: %w", err)
		}
		log.Println("stopping recording")
		if err := rec.Stop(); err != nil {
			return fmt.Errorf("stop recording: %w", err)
		}
		var wav bytes.Buffer
		if err := rec.WriteWAV(&wav); err != nil {
			return fmt.Errorf("write wav: %w", err)
		}
		if d.Dump {
			f, err := os.CreateTemp("", "whisperd-*.wav")
			if err != nil {
				return fmt.Errorf("dump wav: %w", err)
			}
			if _, err := wav.WriteTo(f); err != nil {
				f.Close()
				return fmt.Errorf("dump wav: %w", err)
			}
			if err := f.Close(); err != nil {
				return fmt.Errorf("dump wav: %w", err)
			}
			log.Printf("dumped: %s", f.Name())
		}
		tray.SetStatus(tray.Transcribing)
		log.Println("transcribing ...")
		text, err := d.Client.Transcribe(ctx, &wav)
		if err != nil {
			return fmt.Errorf("transcribe: %w", err)
		}
		log.Printf("emitting: %s", text)
		if err := uinput.EmitText(d.Output, text); err != nil {
			return fmt.Errorf("emit text: %w", err)
		}
	}
}
