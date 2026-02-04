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

### Key Codes

For available key codes to use with the `-key` flag, see [internal/inputcodes/codes.go](internal/inputcodes/codes.go).

## Usage

1. Set your OpenAI API key:
   ```sh
   export OPENAI_API_KEY="your-api-key-here"
   ```

2. Find your input device:
   ```sh
   ls /dev/input/event*
   # or use evtest to identify the correct device
   sudo evtest
   ```

3. Build and run the daemon:
   ```sh
   go build -o whisperd .
   sudo ./whisperd -input /dev/input/event3
   ```

4. Hold the configured hotkey to dictate text.

## Permissions

- Add your user to the `input` group:
  ```sh
  sudo usermod -aG input $USER
  ```
- Or set udev rules to allow access to `/dev/uinput` and input devices without sudo. 

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