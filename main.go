package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/icholy/whisperd/internal/daemon"
	"github.com/icholy/whisperd/internal/inputcodes"
	"github.com/icholy/whisperd/internal/openai"
	"github.com/icholy/whisperd/internal/tray"
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
	d := &daemon.Daemon{
		Log:     slog.Default(),
		Input:   input,
		Output:  output,
		Client:  openai.Client{APIKey: openaiKey, BaseURL: openaiBaseURL},
		KeyCode: uint16(keyCode),
		Dump:    dump,
	}
	ctx := context.Background()
	tray.Run(func() {
		go func() {
			if err := d.Run(ctx); err != nil {
				log.Fatal(err)
			}
		}()
	}, nil)
}
