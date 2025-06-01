package main

import (
	"bytes"
	"context"
	"log"
	"os"

	"github.com/icholy/whisperd/internal/evdev"
	"github.com/icholy/whisperd/internal/inputcodes"
	"github.com/icholy/whisperd/internal/openapi"
	"github.com/icholy/whisperd/internal/pipewire"
	"github.com/icholy/whisperd/internal/uinput"
)

func main() {
	// open input keyboard
	input, err := os.Open("/dev/input/event5")
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()
	// create output keyboard
	output, err := uinput.Create("whisperd")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()
	defer uinput.Destroy(output)
	// setup openai
	client := openai.Client{
		APIKey: os.Getenv("OPENAI_API_KEY"),
	}

	ctx := context.Background()

	for {
		log.Println("waiting for key")
		recctx, err := evdev.KeyDownContext(input, inputcodes.KEY_MAIL)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("starting recording")
		rec, err := pipewire.Record(recctx, pipewire.Options{
			SampleRate:  16000,
			NumChannels: 1,
		})
		if err != nil {
			log.Fatal(err)
		}
		<-recctx.Done()
		log.Println("stopping recording")
		if err := rec.Stop(); err != nil {
			log.Fatal(err)
		}
		var wav bytes.Buffer
		if err := rec.WriteWAV(&wav); err != nil {
			log.Fatal()
		}
		log.Println("transcribing")
		text, err := client.Transcribe(ctx, &wav, "audio.wav")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("emitting: %s", text)
		if err := uinput.EmitText(output, text); err != nil {
			log.Fatal(err)
		}
	}
}
