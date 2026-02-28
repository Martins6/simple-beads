#!/bin/bash

# Simple Beads (sb) Installation Script
# For Unix systems only (macOS and Linux)
# No sudo required - installs to ~/.local/bin by default

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_success() { echo -e "${GREEN}✓${NC} $1"; }
print_error() { echo -e "${RED}✗${NC} $1"; }
print_info() { echo -e "${YELLOW}ℹ${NC} $1"; }
print_cmd() { echo -e "${BLUE}$1${NC}"; }

# Check if Unix system
case "$(uname -s)" in
    Linux|Darwin) ;;  # Supported
    *) 
        print_error "Unsupported operating system: $(uname -s)"
        print_error "sbeads only supports Unix systems (macOS and Linux)"
        exit 1
        ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) 
        print_error "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Installation directories
LOCAL_DIR="$HOME/.local/bin"
GO_BIN="$(go env GOPATH 2>/dev/null)/bin" || "$HOME/go/bin"
BINARY_NAME="sb"
REPO_URL="https://github.com/user/sbeads"

# Determine install location
get_install_dir() {
    # Prefer ~/.local/bin (no sudo needed)
    if [ -d "$LOCAL_DIR" ] || mkdir -p "$LOCAL_DIR" 2>/dev/null; then
        echo "$LOCAL_DIR"
        return 0
    fi
    
    # Fall back to Go bin directory
    if [ -d "$GO_BIN" ]; then
        echo "$GO_BIN"
        return 0
    fi
    
    # Try to create Go bin directory
    if mkdir -p "$GO_BIN" 2>/dev/null; then
        echo "$GO_BIN"
        return 0
    fi
    
    # Last resort: check if /usr/local/bin is writable
    if [ -w /usr/local/bin ]; then
        echo "/usr/local/bin"
        return 0
    fi
    
    return 1
}

# Main install function
main() {
    print_info "Installing Simple Beads (sb)"
    print_info "============================"
    print_info "Platform: $(uname -s)/$ARCH"
    echo ""

    # Determine install location
    TARGET_DIR=$(get_install_dir)
    if [ $? -ne 0 ]; then
        print_error "Could not find a suitable installation directory."
        echo ""
        echo "Please create one of these directories manually:"
        print_cmd "  mkdir -p ~/.local/bin"
        echo "or"
        print_cmd "  mkdir -p ~/go/bin"
        echo ""
        echo "Then add it to your PATH in ~/.bashrc or ~/.zshrc:"
        print_cmd "  export PATH=\"\$PATH:\$HOME/.local/bin\""
        exit 1
    fi

    # Check if we're in source directory
    if [ -f "go.mod" ] && [ -d "cmd" ]; then
        print_info "Building from source..."
        
        # Check for Go
        if ! command -v go &> /dev/null; then
            print_error "Go is not installed"
            print_info "Install Go from: https://golang.org/doc/install"
            exit 1
        fi
        
        # Build
        go build -o "$BINARY_NAME" .
        print_success "Built $BINARY_NAME"
    else
        # Download pre-built binary
        print_info "Downloading pre-built binary..."
        
        OS=$(uname -s | tr '[:upper:]' '[:lower:]')
        BINARY_URL="$REPO_URL/releases/latest/download/sb-${OS}-${ARCH}"
        
        if command -v curl &> /dev/null; then
            curl -fsSL "$BINARY_URL" -o "$BINARY_NAME" || {
                print_error "Failed to download binary"
                print_info "You can build from source instead:"
                print_cmd "  git clone $REPO_URL"
                print_cmd "  cd sbeads"
                print_cmd "  go build -o sb ."
                exit 1
            }
        elif command -v wget &> /dev/null; then
            wget -q "$BINARY_URL" -O "$BINARY_NAME" || {
                print_error "Failed to download binary"
                print_info "You can build from source instead:"
                print_cmd "  git clone $REPO_URL"
                print_cmd "  cd sbeads"
                print_cmd "  go build -o sb ."
                exit 1
            }
        else
            print_error "Need curl or wget to download"
            exit 1
        fi
        
        chmod +x "$BINARY_NAME"
        print_success "Downloaded $BINARY_NAME"
    fi

    # Install
    print_info "Installing to $TARGET_DIR..."
    cp "$BINARY_NAME" "$TARGET_DIR/"
    chmod +x "$TARGET_DIR/$BINARY_NAME"
    print_success "Installed to $TARGET_DIR/$BINARY_NAME"

    # Cleanup
    rm -f "$BINARY_NAME"

    # Verify
    if [ -x "$TARGET_DIR/$BINARY_NAME" ]; then
        print_success "Installation complete!"
        echo ""
        echo "Quick start:"
        print_cmd "  sb init              # Initialize sbeads"
        print_cmd "  sb create \"My task\"  # Create first task"
        print_cmd "  sb list              # List tasks"
        print_cmd "  sb --help            # Show all commands"
        echo ""
        
        # Check if in PATH
        if [[ ":$PATH:" != *":$TARGET_DIR:"* ]]; then
            print_info "Add $TARGET_DIR to your PATH:"
            print_cmd "  export PATH=\"\$PATH:$TARGET_DIR\""
            echo ""
            print_info "Add this to your ~/.bashrc or ~/.zshrc to make it permanent."
        fi
    else
        print_error "Installation verification failed"
        exit 1
    fi
}

main "$@"
