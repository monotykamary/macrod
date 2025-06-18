#!/bin/bash

echo "Testing hotkey functionality..."
echo "1. Make sure the daemon is running"
echo "2. Enable a macro with hotkey 'ctrl+shift+1' in the TUI"
echo "3. Press Ctrl+Shift+1 to trigger the macro"
echo ""
echo "The daemon should log when hotkeys are registered and triggered."
echo ""
echo "Current daemon process:"
ps aux | grep macrod-daemon | grep -v grep