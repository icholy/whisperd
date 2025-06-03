package main

import (
	"bytes"
	"context"
	"flag"
	"log"
	"os"

	"github.com/icholy/whisperd/internal/evdev"
	"github.com/icholy/whisperd/internal/inputcodes"
	"github.com/icholy/whisperd/internal/openai"
	"github.com/icholy/whisperd/internal/pipewire"
	"github.com/icholy/whisperd/internal/uinput"
)

func main() {
	var inputPath, openaiKey string
	var keyCode int
	flag.StringVar(&inputPath, "input", "", "device path to use. Ex: /dev/input/eventX")
	flag.IntVar(&keyCode, "key", int(inputcodes.KEY_MAIL), "Key code to use")
	flag.StringVar(&openaiKey, "openai.key", "", "OpenAI API Key")
	flag.Parse()
	if openaiKey == "" {
		openaiKey = os.Getenv("OPENAI_API_KEY")
	}
	if openaiKey == "" {
		log.Fatal("no api key found")
	}
	if inputPath == "" {
		log.Fatal("missing input device path")
	}
	// open input keyboard
	input, err := os.Open(inputPath)
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
	client := openai.Client{APIKey: openaiKey}
	ctx := context.Background()
	for {
		log.Println("waiting for key down")
		if err := evdev.WaitForKey(input, uint16(keyCode), 1); err != nil {
			log.Fatal(err)
		}
		log.Println("starting recording")
		rec, err := pipewire.Record(ctx, pipewire.Options{
			SampleRate:  16000,
			NumChannels: 1,
		})
		log.Println("waiting for key up")
		if err := evdev.WaitForKey(input, uint16(keyCode), 0); err != nil {
			log.Fatal(err)
		}
		if err != nil {
			log.Fatal(err)
		}
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
