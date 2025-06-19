# macrod - Macro Daemon for macOS

A powerful TUI-based macro recording and playback daemon for macOS, designed for automating gameplay combos in RetroArch and other applications.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/Platform-macOS-000000?style=flat&logo=apple)
![License](https://img.shields.io/badge/License-MIT-blue?style=flat)

## Installation

### Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/monotykamary/macrod/main/install.sh | bash
```

Or if you prefer to review the script first:

```bash
curl -fsSL https://raw.githubusercontent.com/monotykamary/macrod/main/install.sh -o install.sh
cat install.sh  # Review the script
bash install.sh
```

### Manual Installation

1. Download the latest release for your architecture from the [releases page](https://github.com/monotykamary/macrod/releases)
2. Extract the archive: `tar -xzf macrod-darwin-*.tar.gz`
3. Move the binaries to `/usr/local/bin/`:
   ```bash
   sudo mv macrod-daemon macrod-tui /usr/local/bin/
   ```
4. Make them executable: `chmod +x /usr/local/bin/macrod-*`

### Building from Source

```bash
git clone https://github.com/monotykamary/macrod.git
cd macrod
make build
# Binaries will be in ./bin/
```

## Features

- 🎮 **Game-Ready**: Designed specifically for gameplay automation and combo execution
- 🎨 **Beautiful TUI**: Interactive terminal interface built with Bubble Tea
- ⏺️ **Global Key Recording**: Native macOS CGEventTap for accurate capture
- 📝 **Macro Management**: Full CRUD operations with edit functionality
- ⌨️ **Hotkey Triggers**: Press custom hotkeys to instantly execute macros
- 🔄 **Daemon Architecture**: Background service for system-wide functionality
- 💾 **Persistent Storage**: JSON-based storage for macro persistence
- 🛡️ **Safety Features**: Confirmation dialogs and smart recording controls

## Quick Start

After installation, you'll have the `macrod` command available:

```bash
# Start the daemon in the background
macrod start

# Launch the TUI interface
macrod

# Check daemon status
macrod status

# Stop the daemon
macrod stop
```

### First Time Setup

1. **Grant Accessibility Permissions**: When you first run `macrod start`, macOS will prompt you to grant accessibility permissions. This is required for global key capture.

2. **Start Recording**: In the TUI, press `r` to start recording a macro. Press keys for your combo, then press `Esc` to finish.

3. **Set a Hotkey**: After recording, set a hotkey (like `ctrl+1`) to trigger your macro.

4. **Play Your Macro**: Press `p` in the TUI or use your hotkey anywhere on your system!

## Architecture

The project uses a client-server architecture:

1. **Daemon** (`macrod-daemon`): Background service that handles:
   - Global key capture (with user permission)
   - Macro playback execution
   - Hotkey registration
   - Persistent storage management

2. **TUI Client** (`macrod-tui`): Terminal interface that provides:
   - Visual macro management
   - Recording interface
   - Real-time daemon status
   - Interactive editing

Communication happens via Unix sockets at `/tmp/macrod.sock`.
make run-daemon
```

2. **In another terminal, launch the TUI**:
```bash
make run-tui
```

3. **Record your first macro**:
   - Press `r` to start recording
   - Type your key sequence (e.g., "Hello World!")
   - Press `Esc` to stop recording
   - Fill in the macro details
   - Press `Enter` to save

4. **Use your macro**:
   - Enable it with `Space`
   - Press the hotkey anywhere (e.g., Ctrl+Shift+1)
   - Or press `p` in the TUI to play manually

## Usage

### TUI Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate through macros |
| `←/→` or `h/l` | Scroll horizontally |
| `Space` | Toggle macro enable/disable |
| `e` | Edit selected macro |
| `d` | Delete selected macro (with confirmation) |
| `r` | Record new macro |
| `p` | Play selected macro |
| `?` | Toggle help |
| `q` | Quit |

### Recording Mode

When recording (`r` key):
1. **Global Capture**: All keystrokes are captured system-wide
2. **Visual Feedback**: Recording status shown in real-time
3. **Smart Stop**: Press `Esc` to stop recording (Esc not included in macro)
4. **Form Entry**: Fill in name, description, and hotkey (not recorded)
5. **Save**: Press `Enter` on the hotkey field to save

**Note**: The TUI disables during recording to prevent interference. All keys are captured globally by the daemon.

### macOS Permissions

The daemon requires Accessibility permissions for global key capture:

1. Open **System Preferences** → **Security & Privacy** → **Privacy** → **Accessibility**
2. Click the lock to make changes
3. Add `macrod-daemon` to the list (drag from Finder or use `+` button)
4. Ensure the checkbox is checked

## Features in Detail

### Macro Recording
- Native CGEventTap for system-wide capture
- Smart key detection (capitals → lowercase + shift)
- Special character support (!, @, # → base key + shift)
- Precise timing preservation between keystrokes
- Automatic Esc key filtering (stop trigger not recorded)

### Macro Playback
- Native CGEventPost for accurate reproduction
- Configurable minimum delay (50ms default)
- Full modifier key support (Cmd, Ctrl, Alt, Shift)
- Background hotkey triggers
- Manual playback from TUI

### Macro Management
- **Create**: Record new macros with custom metadata
- **Read**: View all macros in a formatted table
- **Update**: Edit macro name, description, and hotkey
- **Delete**: Remove macros with confirmation dialog
- **Toggle**: Enable/disable macros without deletion

### Storage
- Persistent JSON storage at `~/.config/macrod/macros.json`
- Automatic save on changes
- Example macros created on first run

## Project Structure

```
macrod/
├── cmd/
│   ├── daemon/        # Daemon entry point
│   └── tui/           # TUI client entry point
├── internal/
│   ├── ipc/           # Unix socket communication
│   ├── keylogger/     # Key capture and playback
│   ├── macro/         # Macro management logic
│   └── storage/       # Persistence layer
├── pkg/
│   └── models/        # Shared data models
├── Makefile           # Build automation
└── README.md
```

## Development

### Running Tests

```bash
make test
```

### Development Mode

```bash
# Terminal 1: Run daemon with live reload
make dev-daemon

# Terminal 2: Run TUI with live reload
make dev-tui
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Troubleshooting

### "Accessibility permissions required" error
- Ensure `macrod-daemon` is added to Accessibility in System Preferences
- Try removing and re-adding the permission
- Restart the daemon after granting permissions

### "Daemon not running" in TUI
- Check if daemon is running: `ps aux | grep macrod-daemon`
- Check daemon logs for errors
- Ensure no other process is using `/tmp/macrod.sock`

### Recorded keys not playing back
- Verify the macro is enabled (✅ in status column)
- Check if the target application is in focus
- Some applications may have anti-automation measures

## Roadmap

### Completed ✅
- [x] Basic macro recording and playback
- [x] TUI with table view
- [x] Daemon/client architecture
- [x] Persistent storage
- [x] Edit functionality
- [x] Delete confirmation
- [x] Native CGEventTap implementation for key capture
- [x] Global hotkey triggers
- [x] Smart key recording (capitals, special chars)
- [x] Pause/resume recording controls

### Future Enhancements 🚀
- [ ] Macro groups/categories
- [ ] Import/export functionality
- [ ] Cross-platform support (Linux, Windows)
- [ ] Macro scripting language
- [ ] Performance optimizations (see below)
- [ ] Mouse event recording
- [ ] Conditional macros
- [ ] Loop/repeat functionality

## Performance Notes

### Current Implementation
The macro playback currently uses a minimum 50ms delay between keystrokes for reliability. This ensures compatibility across different applications but may feel slow for gaming combos.

### Optimization Opportunities

1. **Adjustable Playback Speed**
   - Add per-macro speed multiplier (0.1x - 10x)
   - Allow microsecond precision for gaming
   - Profile-based timing for different applications

2. **Batch Key Events**
   - Send multiple simultaneous keys as single event
   - Optimize modifier key handling
   - Reduce CGEventPost overhead

3. **Direct Input Injection**
   - Bypass event queue for faster delivery
   - Use IOKit for lower-level access
   - Application-specific injection methods

4. **Recording Optimization**
   - Filter duplicate modifier events
   - Compress timing data
   - Smart event coalescing

### Configuration Ideas
```json
{
  "playback": {
    "defaultDelay": 50,
    "minDelay": 0,
    "speedMultiplier": 1.0,
    "gamingMode": false,
    "batchEvents": true
  }
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - The excellent TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - For beautiful terminal styling
- [keybd_event](https://github.com/micmonay/keybd_event) - Cross-platform keyboard simulation

---

Made with ❤️ for gamers and automation enthusiasts