#!/bin/bash

# Test script for macrod

echo "🎮 Macro Daemon Test Script"
echo "=========================="
echo ""

# Check if binaries exist
if [ ! -f "bin/macrod-daemon" ] || [ ! -f "bin/macrod-tui" ]; then
    echo "❌ Binaries not found. Building..."
    make build
fi

echo "📝 Instructions:"
echo "1. Run the daemon in one terminal: ./bin/macrod-daemon"
echo "2. Run the TUI in another terminal: ./bin/macrod-tui"
echo ""
echo "TUI Controls:"
echo "- ↑/↓ or j/k: Navigate macros"
echo "- Space/Enter: Toggle enable/disable"
echo "- p: Play selected macro"
echo "- r: Record new macro"
echo "- d: Delete macro"
echo "- ?: Show help"
echo "- q: Quit"
echo ""
echo "Recording Mode:"
echo "- Type any keys to record them"
echo "- Tab: Navigate between name/description/hotkey fields"
echo "- Enter: Save macro (when on hotkey field)"
echo "- Esc: Cancel recording"
echo ""
echo "Note: On macOS, the daemon needs Accessibility permissions to capture global keys."
echo "Currently, the key capture is simulated for testing."