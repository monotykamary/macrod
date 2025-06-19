#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="monotykamary/macrod"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.config/macrod"
TEMP_DIR=$(mktemp -d)

# Cleanup on exit
trap 'rm -rf "$TEMP_DIR"' EXIT

echo -e "${GREEN}macrod installer${NC}"
echo "===================="
echo ""

# Check if running on macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
    echo -e "${RED}Error: macrod currently only supports macOS${NC}"
    exit 1
fi

# Check for required tools
if ! command -v curl &> /dev/null; then
    echo -e "${RED}Error: curl is required but not installed${NC}"
    exit 1
fi

if ! command -v tar &> /dev/null; then
    echo -e "${RED}Error: tar is required but not installed${NC}"
    exit 1
fi

# Detect architecture
ARCH=$(uname -m)
if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "arm64" ]]; then
    ARCH="arm64"
else
    echo -e "${RED}Error: Unsupported architecture: $ARCH${NC}"
    exit 1
fi

echo -e "Detected architecture: ${GREEN}$ARCH${NC}"

# Get latest release
echo -e "\nFetching latest release..."
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [[ -z "$LATEST_RELEASE" ]]; then
    echo -e "${RED}Error: Could not fetch latest release${NC}"
    exit 1
fi

echo -e "Latest version: ${GREEN}$LATEST_RELEASE${NC}"

# Download URL
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/macrod-darwin-$ARCH.tar.gz"

# Download the release
echo -e "\nDownloading macrod..."
if ! curl -L -o "$TEMP_DIR/macrod.tar.gz" "$DOWNLOAD_URL"; then
    echo -e "${RED}Error: Failed to download macrod${NC}"
    exit 1
fi

# Extract the archive
echo -e "Extracting files..."
cd "$TEMP_DIR"
if ! tar -xzf macrod.tar.gz; then
    echo -e "${RED}Error: Failed to extract archive${NC}"
    exit 1
fi

# Check if binaries exist
if [[ ! -f "macrod-daemon" ]] || [[ ! -f "macrod-tui" ]]; then
    echo -e "${RED}Error: Required binaries not found in archive${NC}"
    exit 1
fi

# Make binaries executable
chmod +x macrod-daemon macrod-tui

# Create config directory
echo -e "\nCreating config directory..."
mkdir -p "$CONFIG_DIR"

# Check if we need sudo for installation
if [[ -w "$INSTALL_DIR" ]]; then
    SUDO=""
else
    echo -e "${YELLOW}Note: sudo required to install to $INSTALL_DIR${NC}"
    SUDO="sudo"
fi

# Install binaries
echo -e "\nInstalling binaries to $INSTALL_DIR..."
$SUDO mv macrod-daemon "$INSTALL_DIR/"
$SUDO mv macrod-tui "$INSTALL_DIR/"

# Create convenience script
echo -e "Creating macrod command..."
cat > macrod << 'EOF'
#!/bin/bash

case "$1" in
    daemon)
        exec macrod-daemon
        ;;
    tui|"")
        exec macrod-tui
        ;;
    start)
        echo "Starting macrod daemon in background..."
        macrod-daemon &
        echo "Daemon started with PID $!"
        ;;
    stop)
        echo "Stopping macrod daemon..."
        pkill -f macrod-daemon
        ;;
    status)
        if pgrep -f macrod-daemon > /dev/null; then
            echo "macrod daemon is running"
        else
            echo "macrod daemon is not running"
        fi
        ;;
    help|--help|-h)
        echo "macrod - Macro daemon for macOS"
        echo ""
        echo "Usage:"
        echo "  macrod [command]"
        echo ""
        echo "Commands:"
        echo "  tui     Launch the TUI interface (default)"
        echo "  daemon  Run the daemon in foreground"
        echo "  start   Start the daemon in background"
        echo "  stop    Stop the daemon"
        echo "  status  Check daemon status"
        echo "  help    Show this help message"
        ;;
    *)
        echo "Unknown command: $1"
        echo "Run 'macrod help' for usage information"
        exit 1
        ;;
esac
EOF

chmod +x macrod
$SUDO mv macrod "$INSTALL_DIR/"

# Check for accessibility permissions
echo -e "\n${YELLOW}Important: Accessibility Permissions${NC}"
echo "macrod requires accessibility permissions to capture global keystrokes."
echo "When you first run the daemon, macOS will prompt you to grant permissions."
echo ""
echo "To grant permissions manually:"
echo "1. Open System Preferences > Security & Privacy > Privacy"
echo "2. Select Accessibility from the left sidebar"
echo "3. Click the lock to make changes"
echo "4. Add macrod-daemon to the list"

# Success message
echo -e "\n${GREEN}✅ Installation complete!${NC}"
echo ""
echo "To get started:"
echo "  1. Run 'macrod start' to start the daemon"
echo "  2. Run 'macrod' to open the TUI"
echo ""
echo "For more information, visit: https://github.com/$REPO"