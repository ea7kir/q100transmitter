#!/bin/bash

# updated to GO 1.22 on March 2 2024

GOVERSION=1.22.0

echo Update Pi OS
sudo apt update
sudo apt -y full-upgrade
sudo apt -y autoremove
sudo apt clean

###################################################

echo Installing Go $GOVERSION
GOFILE=go$GOVERSION.linux-arm64.tar.gz
cd /usr/local
sudo rm -rf go
sudo wget https://go.dev/dl/$GOFILE
sudo tar -C /usr/local -xzf $GOFILE
cd /home/pi/Q100/q100transmitter
go mod tidy
go version

