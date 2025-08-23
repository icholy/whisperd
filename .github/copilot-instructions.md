# whisperd - Voice-to-Text Linux Daemon

A Linux daemon for voice-to-text typing using OpenAI Whisper. Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

### Quick Start - Build and Basic Validation
- Bootstrap and build the repository:
  - `go mod download` -- downloads dependencies. Takes <30 seconds. NEVER CANCEL.
  - `go mod tidy` -- resolves dependencies. Takes <5 seconds. 
  - `go build -o whisperd .` -- builds the daemon. Takes <1 second (measured: ~0.3s). NEVER CANCEL, set timeout to 60+ seconds.
- Basic validation that always works:
  - `go fmt ./...` -- format code
  - `go vet ./...` -- static analysis
  - `./whisperd --help` -- verify binary works and shows usage
- Test that build produces working binary:
  - `./whisperd -input /dev/input/event0` should fail with "no api key found" 
  - `OPENAI_API_KEY=test ./whisperd -input /dev/input/event0` should fail with "permission denied" (either input device or uinput access)

### System Requirements for Full Testing
- **Linux only** - uses Linux input devices (/dev/input/eventX) and uinput (/dev/uinput)
- **PipeWire**: `pw-cat` command must be available (`sudo apt-get install pipewire-bin`)
- **Input tools**: `evtest` for device identification (`sudo apt-get install evtest`)
- **Permissions**: Access to `/dev/input/eventX` and `/dev/uinput` (requires root OR proper udev rules)
- **Audio system**: Working PipeWire server for audio recording
- **Go 1.21+**: Runtime requirement (go.mod specifies 1.24.2)

### System Setup for Full Development
```bash
# Install system dependencies (Ubuntu/Debian)
sudo apt-get update
sudo apt-get install -y pipewire-bin evtest

# Check available input devices
ls /dev/input/event*
sudo evtest  # Interactive tool to identify correct device

# Add user to input group (logout/login required)
sudo usermod -aG input $USER

# Test PipeWire audio recording
pw-cat --record --format s16 --rate 16000 --channels 1 - > /tmp/test.wav
# Ctrl+C after a few seconds, check that test.wav was created
```

## Testing and Validation

### What Can Always Be Tested (No Special Permissions)
- **Build process**: `go build -o whisperd .`
- **Code formatting**: `go fmt ./...`
- **Static analysis**: `go vet ./...`
- **Help output**: `./whisperd --help`
- **API key validation**: Error handling when OPENAI_API_KEY is missing
- **Syntax validation**: Go compiler catches syntax errors during build

### What Requires System Permissions
- **Input device access**: Reading from `/dev/input/eventX` 
- **uinput device creation**: Writing to `/dev/uinput` for text emission
- **Audio recording**: PipeWire `pw-cat` functionality
- **Complete workflow**: Hotkey detection → audio recording → transcription → text typing

### Manual Validation Scenarios
After making changes, always test the complete workflow if you have proper permissions:

1. **Setup test environment**:
   ```bash
   export OPENAI_API_KEY="your-actual-api-key"
   # Find your keyboard device (usually event3 or event4)
   sudo evtest
   ```

2. **Run the daemon**:
   ```bash
   # First identify your keyboard input device
   sudo evtest
   # Look for your keyboard device (usually "AT Translated Set 2 keyboard" or similar)
   # Note the /dev/input/eventX path and exit evtest (Ctrl+C)
   
   # Run whisperd with the correct device
   sudo ./whisperd -input /dev/input/event3  # Replace event3 with your device
   ```

3. **Test workflow**:
   - Open a text editor (notepad, terminal, etc.)
   - Focus the text input area
   - Hold the hotkey (default: KEY_MAIL, code 155)
   - Speak something clearly (e.g., "Hello world test")
   - Release the hotkey
   - Wait for transcription to complete (usually 2-5 seconds)
   - Verify transcribed text appears in the text editor
   
   **Note**: This requires a working OpenAI API key with available credits. The transcription will fail gracefully if the API key is invalid or has no credits.

### Timing Expectations
- **Build**: <1 second (measured: ~0.3s). NEVER CANCEL. Set timeout to 60+ seconds for safety.
- **Dependencies**: <30 seconds for `go mod download`. NEVER CANCEL.
- **go mod tidy**: <1 second (measured: ~0.02s).
- **Tests**: No test suite exists (`go test ./...` shows "no test files")
- **Runtime startup**: Immediate (daemon starts and waits for input)

## Architecture and Key Files

### Main Components
- **main.go**: Main daemon loop, handles CLI flags and coordinates all components
- **internal/evdev**: Linux input device handling for hotkey detection
- **internal/pipewire**: Audio recording using PipeWire's `pw-cat` command
- **internal/openai**: OpenAI Whisper API client for speech transcription
- **internal/uinput**: Linux uinput device creation for text emission
- **internal/inputcodes**: Linux input event constants and key codes
- **internal/wav**: WAV audio file format handling

### Key Configuration
- **Hotkey**: Default KEY_MAIL (155), configurable via `-key` flag
- **Input device**: Must specify device path via `-input` flag
- **API key**: Via `-openai.key` flag or `OPENAI_API_KEY` environment variable
- **Audio format**: 16kHz, 1 channel, signed 16-bit PCM (hardcoded in main.go)

### Dependencies
- **golang.org/x/sys**: Only Go dependency for system calls
- **System binaries**: `pw-cat` (PipeWire), `evtest` (optional, for device identification)

## Common Development Tasks

### Building and Running
```bash
# Standard development cycle
go mod tidy
go build -o whisperd .

# Quick syntax check without building
go vet ./...

# Test various error conditions to verify error handling
OPENAI_API_KEY=dummy sudo ./whisperd -input /dev/input/event0  # May fail at system permissions
OPENAI_API_KEY=test ./whisperd -input /dev/input/event0  # "permission denied" 
./whisperd -input /nonexistent  # "no api key found"
OPENAI_API_KEY=test ./whisperd -input /nonexistent  # "no such file or directory"
```

### Key Code Reference
Available key codes are defined in `internal/inputcodes/codes.go`. Common ones:
- `KEY_MAIL = 155` (default)
- `KEY_F1 = 59`, `KEY_F2 = 60`, etc.
- `KEY_LEFTCTRL = 29`, `KEY_RIGHTCTRL = 97`

### Debugging Common Issues
- **"no api key found"**: Set `OPENAI_API_KEY` environment variable
- **"permission denied" on /dev/input**: Run with `sudo` or add user to `input` group
- **"permission denied" on /dev/uinput**: Run with `sudo` or configure udev rules (may occur before input device access)
- **"Host is down" from pw-cat**: PipeWire server not running or audio system not configured
- **"device or resource busy"**: Input device already in use by another process

### Code Style and Linting
- **No custom linters**: Repository uses standard Go tooling only
- **Format code**: Always run `go fmt ./...` before committing
- **Static analysis**: Always run `go vet ./...` before committing
- **No tests**: Repository currently has no test suite

## Limitations in Containerized Environments

### What Works in Containers
- Building the application
- Code formatting and static analysis  
- Syntax validation
- Help output and basic CLI flag testing

### What Doesn't Work in Containers
- Audio recording (no PipeWire server)
- Input device access (no /dev/input devices)
- uinput device creation (no /dev/uinput)
- Complete end-to-end testing

**Note**: This is a system-level daemon that requires a full Linux desktop environment with audio and input device access. Container-based testing is inherently limited for this type of application.

## Validation Checklist

Always validate changes using this checklist:

- [ ] `go mod tidy` completes successfully
- [ ] `go build -o whisperd .` builds without errors (<2s)
- [ ] `go fmt ./...` produces no output (code already formatted)
- [ ] `go vet ./...` produces no warnings
- [ ] `./whisperd --help` displays usage information
- [ ] `./whisperd -input /tmp` fails with appropriate error message
- [ ] If system permissions available: Complete workflow test with real input device
- [ ] If OpenAI API key available: Test actual transcription functionality

Remember: This daemon requires root privileges or specific system configuration. Many validation steps can only be performed on a properly configured Linux desktop system.