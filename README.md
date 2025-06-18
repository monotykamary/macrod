# macrod - Macro Daemon for macOS

A powerful TUI-based macro recording and playback daemon for macOS, designed for automating gameplay combos in RetroArch and other applications.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/Platform-macOS-000000?style=flat&logo=apple)
![License](https://img.shields.io/badge/License-MIT-blue?style=flat)

## Features

- 🎮 **Game-Ready**: Designed specifically for gameplay automation and combo execution
- 🎨 **Beautiful TUI**: Interactive terminal interface built with Bubble Tea
- ⏺️ **Global Key Recording**: Capture key sequences with precise timing
- 📝 **Macro Management**: Full CRUD operations with edit functionality
- ⌨️ **Hotkey Support**: Assign custom hotkeys to trigger macros instantly
- 🔄 **Daemon Architecture**: Background service for system-wide functionality
- 💾 **Persistent Storage**: JSON-based storage for macro persistence
- 🛡️ **Safety Features**: Confirmation dialogs for destructive actions

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

## Installation

### Prerequisites

- macOS 10.15 or later
- Go 1.21 or later
- Accessibility permissions (for global key capture)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/macrod.git
cd macrod

# Build both components
make build

# Or build individually
make build-daemon
make build-tui
```

## Quick Start

1. **Start the daemon** (first time will request accessibility permissions):
```bash
make run-daemon
```

2. **In another terminal, launch the TUI**:
```bash
make run-tui
```

3. **Record your first macro**:
   - Press `r` to start recording
   - Perform your key sequence
   - Press `Esc` to stop recording
   - Fill in the macro details
   - Press `Enter` to save

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
1. **Key Capture Phase**: Press your key sequence, `Esc` to finish
2. **Details Phase**: Fill in name, description, and hotkey
3. **Save**: Press `Enter` on the hotkey field to save

### macOS Permissions

The daemon requires Accessibility permissions for global key capture:

1. Open **System Preferences** → **Security & Privacy** → **Privacy** → **Accessibility**
2. Click the lock to make changes
3. Add `macrod-daemon` to the list (drag from Finder or use `+` button)
4. Ensure the checkbox is checked

## Features in Detail

### Macro Recording
- Captures exact key sequences with timing
- Records modifier keys (Cmd, Ctrl, Alt, Shift)
- Preserves delays between keystrokes
- Visual feedback during recording

### Macro Playback
- Accurate timing reproduction
- Modifier key support
- Background execution via hotkeys
- Manual trigger from TUI

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

- [x] Basic macro recording and playback
- [x] TUI with table view
- [x] Daemon/client architecture
- [x] Persistent storage
- [x] Edit functionality
- [x] Delete confirmation
- [ ] CGEventTap implementation for native key capture
- [ ] Hotkey triggers
- [ ] Macro groups/categories
- [ ] Import/export functionality
- [ ] Cross-platform support (Linux, Windows)
- [ ] Macro scripting language

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - The excellent TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - For beautiful terminal styling
- [keybd_event](https://github.com/micmonay/keybd_event) - Cross-platform keyboard simulation

---

Made with ❤️ for gamers and automation enthusiasts