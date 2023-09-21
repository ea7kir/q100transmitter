#!/bin/bash

echo
echo "-------------------------------"
echo "-- Updateing Pi OS"
echo "-------------------------------"
echo

sudo apt update
sudo apt full-upgrade -y
sudo apt autoremove -y
sudo apt clean

