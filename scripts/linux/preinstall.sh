#!/bin/sh

# Access the environment variable
if [ ! -z "$CLOUDCORE_PSK" ]; then
    echo "$CLOUDCORE_PSK" > /etc/cloudcored/psk
    echo "Successfully loaded pre-shared key"
fi