# whisperd

A Linux daemon for voice-to-text typing using OpenAI Whisper.

## Features
- Hold a hotkey to record audio
- Transcribes speech to text using OpenAI Whisper
- Types the text into the focused window

## Requirements
- Access to `/dev/uinput` and input devices (see Permissions)
- PipeWire (`pw-cat`)
- Go 1.21+

## Usage
1. Build and run the daemon:
   ```sh
   go build -o whisperd .
   sudo ./whisperd
   ```
2. Hold the configured hotkey to dictate text.

## Permissions
- Add your user to the `input` group:
  ```sh
  sudo usermod -aG input $USER
  ```
- Or set udev rules to allow access to `/dev/uinput` and input devices without sudo. 