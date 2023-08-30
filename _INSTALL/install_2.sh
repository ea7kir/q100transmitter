#!/bin/bash

echo
echo "-------------------------------"
echo "-- Installing GIT"
echo "-------------------------------"
echo

sudo apt install git -y

echo
echo "-------------------------------"
echo "-- Installing Go"
echo "-------------------------------"
echo

GOVERSION=go1.21.0.linux-arm64.tar.gz
cd /usr/local
sudo wget https://go.dev/dl/$GOVERSION
# sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf $GOVERSION
cd

echo
echo "-------------------------------"
echo "-- Install gioui dependencies"
echo "-------------------------------"
echo

sudo apt install gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev

echo
echo "-------------------------------"
echo "-- Rebooting in 5 seconds"
echo "--"
echo "-- Then run install_3"
echo "-------------------------------"
echo

sleep 5

sudo reboot
