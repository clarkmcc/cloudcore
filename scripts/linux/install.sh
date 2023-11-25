#!/bin/bash

set -eu

# Function definitions
log() {
    echo "$@" >&2
}

error() {
    log "ERROR: $@"
    exit 1
}

detect_os_and_arch() {
    OS=""
    ARCH=$(uname -m)
    PACKAGE_TYPE=""

    if [ -f /etc/os-release ]; then
        . /etc/os-release
        case "$ID" in
            debian|ubuntu|linuxmint)
                OS="debian"
                PACKAGE_TYPE="deb"
                ;;
            centos|fedora|rhel|rocky|almalinux)
                OS="rhel"
                PACKAGE_TYPE="rpm"
                ;;
            alpine)
                OS="alpine"
                PACKAGE_TYPE="apk"
                ;;
            *)
                OS="unsupported"
                ;;
        esac
    fi

    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        armv7l) ARCH="arm5" ;;
        *) ARCH="unsupported" ;;
    esac

    if [ -z "$OS" ] || [ "$ARCH" = "unsupported" ]; then
        error "Unsupported OS or architecture: $OS, $ARCH"
    fi
}

fetch_latest_version_and_construct_package_url() {
    VERSION=$(curl -L -s -H "Accept: application/vnd.github+json" \
                    -H "X-GitHub-Api-Version: 2022-11-28" \
                    https://api.github.com/repos/clarkmcc/cloudcore/releases/latest | \
                    grep '"tag_name":' | cut -d '"' -f 4)
    PACKAGE_URL="https://github.com/clarkmcc/cloudcore/releases/download/${VERSION}/cloudcored_${VERSION}_linux_${ARCH}.${PACKAGE_TYPE}"
}

check_root_privileges() {
    if [ "$(id -u)" -eq 0 ]; then
        SUDO=""
    elif type sudo >/dev/null 2>&1; then
        SUDO="sudo"
    else
        error "Root privileges required. Please run this script as root or install sudo."
    fi
}

install_package() {
    log "Installing CloudCore version $VERSION for $OS on $ARCH architecture"
    curl -L -o cloudcore_package "${PACKAGE_URL}"

    case "$PACKAGE_TYPE" in
        deb)
            $SUDO dpkg -i cloudcore_package
            ;;
        rpm)
            $SUDO rpm -i cloudcore_package
            ;;
        apk)
            $SUDO apk add --allow-untrusted cloudcore_package
            ;;
        *)
            error "Installation method not supported for package type: $PACKAGE_TYPE"
            ;;
    esac
    rm cloudcore_package
}

main() {
    detect_os_and_arch
    fetch_latest_version_and_construct_package_url
    check_root_privileges
    install_package
    log "CloudCore installation complete!"
}

# Execute main function
main
