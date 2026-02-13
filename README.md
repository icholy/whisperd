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

## Configuration

### Command Line Flags

- `-input` - Device path to use (required). Example: `/dev/input/event3`
- `-key` - Key code to use as hotkey (default: 155, which is KEY_MAIL)
- `-openai.key` - OpenAI API Key (can also be set via `OPENAI_API_KEY` environment variable)
- `-openai.baseurl` - OpenAI Base URL (can be used with locally hosted https://speaches.ai)

### Key Codes

For available key codes to use with the `-key` flag, see [internal/inputcodes/codes.go](internal/inputcodes/codes.go).

## Permissions

Add your user to the `input` group:

```sh
sudo usermod -aG input $USER
```

Log out and back in for the group change to take effect.

## Usage

1. Find your input device:
   ```sh
   ls /dev/input/event*
   # or use evtest to identify the correct device
   sudo evtest
   ```

2. Build and install:
   ```sh
   go install .
   ```

3. Run directly:
   ```sh
   whisperd -input /dev/input/event3 -openai.key "your-key-here"
   ```

4. Hold the configured hotkey to dictate text.

## Systemd User Service

To run whisperd as a user service:

1. Create the service file at `~/.config/systemd/user/whisperd.service`:

```ini
[Unit]
Description=Whisper Daemon - Voice To Text
After=network.target
Wants=network.target

[Service]
ExecStart=%h/go/bin/whisperd -input /dev/input/event3 -openai.key "your-key-here"
Restart=always
RestartSec=5

[Install]
WantedBy=default.target
```

2. Enable and start:

```sh
systemctl --user daemon-reload
systemctl --user enable --now whisperd
```

3. View logs:

```sh
journalctl --user -u whisperd -f
```

## System Tray

whisperd shows a system tray icon (gray=idle, red=recording, yellow=transcribing). For X11 environments that only support XEmbed (e.g. i3bar), use the legacy build tag:

```sh
go build -tags legacy_systray
```

## Local Model

Run an OpenAI compatible API in a docker container: https://speaches.ai/installation

```sh
docker run \
  --rm \
  --detach \
  --publish 8000:8000 \
  --name speaches \
  --volume hf-hub-cache:/home/ubuntu/.cache/huggingface/hub \
  --gpus=all \
  ghcr.io/speaches-ai/speaches:latest-cuda
```

Use the `--openai.baseurl` flag to point at it:

``` sh
whisperd --openai.baseurl http://localhost:8000/v1 ...
```