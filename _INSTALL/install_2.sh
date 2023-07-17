#!/bin/bash

echo
echo "-------------------------------"
echo "-- Installing GIT"
echo "-------------------------------"
echo

sudo apt install git -y

git config --global user.name "ea7kir"
git config --global user.email "mikenaylorspain@icloud.com"
git config --global init.defaultBranch main

echo
echo "-------------------------------"
echo "-- Installing Go"
echo
echo "-- this will take some time..."
echo "-------------------------------"
echo

GOVERSION=go1.20.6.linux-arm64.tar.gz
cd /usr/local
sudo wget https://go.dev/dl/$GOVERSION
sudo rm -rf /usr/local/go && tar -C /usr/local -xzf $GOVERSION

sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.20.6.linux-arm64.tar.gz
 
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
