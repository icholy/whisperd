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
	var inputPath, openaiKey, openaiBaseURL string
	var keyCode int
	var dump bool
	flag.StringVar(&inputPath, "input", "", "device path to use. Ex: /dev/input/eventX")
	flag.IntVar(&keyCode, "key", int(inputcodes.KEY_MAIL), "Key code to use")
	flag.StringVar(&openaiKey, "openai.key", "", "OpenAI API Key")
	flag.StringVar(&openaiBaseURL, "openai.baseurl", "", "OpenAI base url")
	flag.BoolVar(&dump, "dump", false, "dump wav contents to files for debugging")
	flag.Parse()
	if openaiKey == "" {
		openaiKey = os.Getenv("OPENAI_API_KEY")
	}
	if openaiKey == "" && openaiBaseURL == "" {
		log.Fatal("no api key found")
	}
	if inputPath == "" {
		log.Fatal("missing input device path")
	}
	// open input keyboard
	input, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("failed to open input device %s: %v", inputPath, err)
	}
	defer input.Close()
	// create output keyboard
	output, err := uinput.Create("whisperd")
	if err != nil {
		log.Fatalf("failed to create uinput device: %v", err)
	}
	defer output.Close()
	defer uinput.Destroy(output)
	// setup openai
	client := openai.Client{APIKey: openaiKey, BaseURL: openaiBaseURL}
	ctx := context.Background()
	for {
		log.Println("waiting for key down")
		if err := evdev.WaitForKey(input, uint16(keyCode), 1); err != nil {
			log.Fatalf("failed to wait for key down: %v", err)
		}
		log.Println("starting recording")
		rec, err := pipewire.Record(ctx, pipewire.Options{
			SampleRate:  16000,
			NumChannels: 1,
		})
		if err != nil {
			log.Fatalf("failed to start recording: %v", err)
		}
		log.Println("waiting for key up")
		if err := evdev.WaitForKey(input, uint16(keyCode), 0); err != nil {
			log.Fatalf("failed to wait for key up: %v", err)
		}
		log.Println("stopping recording")
		if err := rec.Stop(); err != nil {
			log.Fatalf("failed to stop recording: %v", err)
		}
		var wav bytes.Buffer
		if err := rec.WriteWAV(&wav); err != nil {
			log.Fatalf("failed to write WAV data: %v", err)
		}
		if dump {
			f, err := os.CreateTemp("", "whisperd-*.wav")
			if err != nil {
				log.Fatalf("failed to dump wav: %v", err)
			}
			if _, err := wav.WriteTo(f); err != nil {
				log.Fatalf("failed to dump wav: %v", err)
			}
			if err := f.Close(); err != nil {
				log.Fatalf("failed to dump wav: %v", err)
			}
			log.Printf("dumped: %s", f.Name())
		}
		log.Println("transcribing ...")
		text, err := client.Transcribe(ctx, &wav)
		if err != nil {
			log.Fatalf("failed to transcribe audio: %v", err)
		}
		log.Printf("emitting: %s", text)
		if err := uinput.EmitText(output, text); err != nil {
			log.Fatalf("failed to emit text: %v", err)
		}
	}
}
