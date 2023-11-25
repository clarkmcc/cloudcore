#!/bin/bash

CAN_ROOT=''
SUDO=''

if [ "$(id -u)" = 0 ]; then
        CAN_ROOT=1
        SUDO=""
    elif type sudo >/dev/null; then
        CAN_ROOT=1
        SUDO="sudo"
    elif type doas >/dev/null; then
        CAN_ROOT=1
        SUDO="doas"
fi

if [ "$CAN_ROOT" != "1" ]; then
        echo "could not obtain root or sudo access, aborting install. re-run this script as root or setup sudo"
        return
fi

# Function to fetch the latest version number from GitHub
get_latest_version() {
    curl -L -s \
    -H "Accept: application/vnd.github+json" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    https://api.github.com/repos/clarkmcc/cloudcore/releases/latest | \
    grep '"tag_name":' | \
    cut -d '"' -f 4
}

# Function to determine OS
get_os() {
    if [ -f /etc/debian_version ]; then
        echo "debian"
    elif [ -f /etc/alpine-release ]; then
        echo "alpine"
    elif [ -f /etc/arch-release ]; then
        echo "arch"
    elif [ -f /etc/redhat-release ]; then
        echo "rhel"
    else
        echo "unsupported"
    fi
}

# Function to determine architecture using lscpu
get_architecture() {
    lscpu | grep Architecture | cut -d ':' -f 2 | sed 's/ //g'
}

# Function to download and install the package
install_package() {
    local version=$1
    local os=$2
    local arch=$3
    local base_url="https://github.com/clarkmcc/cloudcore/releases/download/${version}"
    local package_name=""
    local installer_command=""
    local package_extension=""

    echo "Installing version $version for $os on $arch architecture"

    # Determine the package name, extension, and installer command based on OS and architecture
    if [ "$os" = "debian" ]; then
        package_extension=".deb"
        installer_command="${SUDO} dpkg -i"
        case $arch in
            x86_64) package_name="cloudcored_${version}_linux_amd64${package_extension}" ;;
            aarch64) package_name="cloudcored_${version}_linux_arm64${package_extension}" ;;
            armv7l) package_name="cloudcored_${version}_linux_arm5${package_extension}" ;;
        esac
    elif [ "$os" = "rhel" ] && [ "$arch" = "aarch64" ]; then
        package_extension=".rpm"
        installer_command="${SUDO} rpm -i"
        package_name="cloudcored_${version}_linux_aarch64${package_extension}"
    else
        echo "Unsupported OS or architecture: $os, $arch"
        return
    fi

    # Construct the download URL
    local url="${base_url}/${package_name}"

    # Download and install the package
    if [ -n "$package_name" ]; then
        echo -n "Fetching $url: "
        curl -L -o "$package_name" "$url"
        echo "done"
        echo "Installing package"
        $installer_command "$package_name"
        rm "$package_name"
    fi
}

# Main installation process
version=$(get_latest_version)
os=$(get_os)
arch=$(get_architecture)

# Parse command-line arguments for the --psk parameter
while [ $# -gt 0 ]; do
    case "$1" in
        --psk)
            CLOUDCORE_PSK="$2"
            shift 2
            ;;
        *)
            break
            ;;
    esac
done


if [ "$os" = "unsupported" ] || [ "$arch" = "unsupported" ]; then
    echo "Error: Unsupported OS or architecture."
    exit 1
fi

export CLOUDCORE_PSK
install_package "$version" "$os" "$arch"
unset CLOUDCORE_PSK
