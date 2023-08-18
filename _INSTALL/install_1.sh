#!/bin/bash

# TODO: ask for call sign

echo
echo "-------------------------------"
echo "-- Updateing the OS"
echo "-------------------------------"
echo

sudo apt update
sudo apt full-upgrade -y
sudo apt autoremove -y
sudo apt clean

echo
echo "-------------------------------"
echo "-- Updating eeprom firmware"
echo "-------------------------------"
echo

sudo rpi-eeprom-update -a

echo
echo "-------------------------------"
echo "-- Updating eeprom firmware"
echo "-------------------------------"
echo

# NOTE: only if advised to do so
# sudo rpi-update

echo
echo "-------------------------------"
echo "-- running rfkill"
echo "-------------------------------"
echo

rfkill block 0
rfkill block 1

echo
echo "-------------------------------"
echo "-- Setting .profile"
echo "-------------------------------"
echo

echo -e '\n\nexport PATH=$PATH:/usr/local/go/bin\n\n' >> /home/pi/.profile

echo
echo "-------------------------------"
echo "-- create Q100 directory"
echo "-------------------------------"
echo

mkdir /home/pi/Q100

echo
echo "-------------------------------"
echo "-- Rebooting in 5 seconds"
echo "--"
echo "-- Then run install_2"
echo "-------------------------------"
echo

sleep 5

sudo reboot
