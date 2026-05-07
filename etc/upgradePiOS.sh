#!/bin/bash

echo "
###################################################
Update Pi OS
###################################################
"

echo Upgradw Pi OS
sudo apt update
sudo apt -y upgrade
sudo apt -y autoremove
sudo apt clean
