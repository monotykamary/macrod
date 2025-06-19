#!/bin/bash
# This is the script users will curl | bash
# It downloads and runs the full installer

set -e

echo "Downloading macrod installer..."
curl -fsSL https://raw.githubusercontent.com/monotykamary/macrod/main/install.sh | bash