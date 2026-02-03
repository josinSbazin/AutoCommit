#!/bin/bash
set -e

REPO="josinSbazin/AutoCommit"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        *)       echo "unsupported" ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64)  echo "amd64" ;;
        aarch64) echo "arm64" ;;
        arm64)   echo "arm64" ;;
        *)       echo "unsupported" ;;
    esac
}

main() {
    OS=$(detect_os)
    ARCH=$(detect_arch)

    if [ "$OS" = "unsupported" ] || [ "$ARCH" = "unsupported" ]; then
        echo "Unsupported OS or architecture"
        exit 1
    fi

    echo "Detected: $OS/$ARCH"

    LATEST=$(curl -sL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$LATEST" ]; then
        echo "Failed to get latest version"
        exit 1
    fi

    echo "Latest version: $LATEST"

    FILENAME="autocommit_${OS}_${ARCH}.tar.gz"
    URL="https://github.com/$REPO/releases/download/$LATEST/$FILENAME"

    echo "Downloading $URL..."

    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT

    curl -sL "$URL" -o "$TMP_DIR/$FILENAME"
    tar -xzf "$TMP_DIR/$FILENAME" -C "$TMP_DIR"

    mkdir -p "$INSTALL_DIR"
    mv "$TMP_DIR/autocommit" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/autocommit"

    echo "Installed to $INSTALL_DIR/autocommit"

    if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
        echo ""
        echo "Add to PATH:"
        echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
    fi

    echo "Done! Run 'autocommit --help' to get started."
}

main
