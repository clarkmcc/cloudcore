#!/bin/bash

# Function to fetch the latest version number from GitHub
get_latest_version() {
    curl -L -s \
    -H "Accept: application/vnd.github+json" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    https://api.github.com/repos/clarkmcc/cloudcore/releases/latest | jq -r .tag_name
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
    local url=""

    echo "Installing version $version for $os on $arch architecture"

    if [ "$os" = "debian" ] && [ "$arch" = "x86_64" ]; then
        url="${base_url}/cloudcored_${version}_linux_amd64.deb"
    elif [ "$os" = "debian" ] && [ "$arch" = "aarch64" ]; then
        url="${base_url}/cloudcored_${version}_linux_arm64.deb"
    elif [ "$os" = "debian" ] && [ "$arch" = "armv7l" ]; then
        url="${base_url}/cloudcored_${version}_linux_arm5.deb"
    # Add more cases for different OS and architecture combinations
    else
        echo "Unsupported OS or architecture: $os, $arch"
        return
    fi

    echo -n "Fetching $url: "
    curl -L -o package.deb "$url"
    echo "done"

    echo "Installing package"
    sudo dpkg -i package.deb
    rm package.deb
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
