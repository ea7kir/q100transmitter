#!/bin/bash

echo "
###################################################
Update Pi OS
###################################################
"

echo Update Pi OS
sudo apt update
sudo apt -y full-upgrade
sudo apt -y autoremove
sudo apt clean
