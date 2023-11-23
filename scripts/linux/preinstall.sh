#!/bin/sh

# Access the environment variable
if [ ! -z "$CLOUDCORE_PSK" ]; then
    mkdir -p /etc/cloudcored
    echo "$CLOUDCORE_PSK" > /etc/cloudcored/psk
    echo "Successfully loaded pre-shared key"
fi