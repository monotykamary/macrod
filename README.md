# macrod - Macro Daemon for macOS

A TUI-based macro recording and playback daemon for automating gameplay combos in RetroArch and other applications.

## Features

- **TUI Interface**: Beautiful terminal interface built with Bubble Tea
- **Macro Recording**: Record key sequences with timing information
- **Macro Management**: Create, edit, enable/disable, and delete macros
- **Hotkey Support**: Assign hotkeys to trigger macros
- **Daemon Architecture**: Background daemon for global key capture
- **Cross-platform Design**: Built for macOS with cross-platform support in mind

## Architecture

The project consists of two main components:

1. **Daemon** (`macrod-daemon`): Runs in the background, captures global keystrokes, and executes macros
2. **TUI Client** (`macrod-tui`): Terminal UI for managing macros and interacting with the daemon

Communication between the TUI and daemon happens via Unix sockets.

## Building

```bash
# Build both components
make build

# Or build individually
make build-daemon
make build-tui
```

## Running

1. Start the daemon (requires accessibility permissions on macOS):
```bash
make run-daemon
# or
./bin/macrod-daemon
```

2. In another terminal, start the TUI:
```bash
make run-tui
# or
./bin/macrod-tui
```

## Usage

### TUI Controls

- `↑/↓` or `j/k`: Navigate through macros
- `Space/Enter`: Toggle macro enable/disable
- `d`: Delete selected macro
- `n`: Create new macro
- `r`: Start/stop recording
- `?`: Show help
- `q`: Quit

### macOS Permissions

The daemon requires Accessibility permissions to capture global keystrokes:

1. Go to System Preferences → Security & Privacy → Privacy → Accessibility
2. Add and enable the `macrod-daemon` executable

## Current Status

The basic architecture is in place with:
- ✅ Project structure and build system
- ✅ TUI with table view for macros
- ✅ Daemon with IPC support
- ✅ Basic storage system
- ⚠️  Key capture implementation (placeholder - needs platform-specific code)
- ⚠️  Key playback implementation (placeholder - needs platform-specific code)

## TODO

- [ ] Implement actual key capture for macOS using CGEventTap
- [ ] Implement key playback functionality
- [ ] Add confirmation dialogs for destructive actions
- [ ] Implement macro editing UI
- [ ] Add macro import/export functionality
- [ ] Cross-platform support (Linux, Windows)

## Development

This project uses:
- Go 1.24+
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI
- [Lipgloss](https://github.com/charmbracelet/lipgloss) for styling

## License

MIT