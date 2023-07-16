#!/bin/bash

cd

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
echo "--"
echo "-- this will take some time..."
echo "-------------------------------"
echo

sudo wget https://go.dev/dl/go1.20.6.linux-arm64.tar.gz
sudo mv go1.20.6.linux-arm64.tar.gz /usr/local/
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.20.6.linux-arm64.tar.gz

echo
echo "-------------------------------"
echo "-- Setting .profile"
echo "-------------------------------"
echo

echo "TODO"

echo
echo "-------------------------------"
echo "-- Install gio"
echo "-------------------------------"
echo

export PATH=$PATH:/usr/local/go/bin
go install gioui.org/cmd/gogio@latest

echo
echo "-------------------------------"
echo "-- Rebooting in 5 seconds"
echo "--"
echo "-- Then run install_3"
echo "-------------------------------"
echo

sleep 5

sudo reboot
