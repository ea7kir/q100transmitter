#!/bin/bash

# Install Q100 Transmitter on txtouch.local
# Orignal design by Michael, EA7KIR

# CONFIFIGURATION
GOVERSION=1.22.5
GIOUIVERSION=7.1
LAN=eth0
PLUTO=eth1
ENCODER=eth2

# DEVICE         TYPE      STATE                                  CONNECTION         
# eth0           ethernet  connected                              Wired connection 1 
# eth1           ethernet  connected                              Wired connection 2 
# lo             loopback  connected (externally)                 lo                 
# eth2           ethernet  connecting (getting IP configuration)  Wired connection 3 

whoami | grep -q pi
if [ $? != 0 ]; then
  echo Install must be performed as user pi
  exit
fi

hostname | grep -q txtouch
if [ $? != 0 ]; then
  echo Install must be performed on host txtouch
  exit
fi

read -p "Enter your callsign in uppercase and press enter " callsign
echo $callsign > /home/pi/Q100/callsign
echo saved to Q100/callsign

while true; do
    read -p "Install q100transmitter using Go $GOVERSION and GIO $GIOUIVERSION (y/n)? " answer
    case ${answer:0:1} in
        y|Y ) break;;
        n|N ) exit;;
        * ) echo "Please answer yes or no.";;
    esac
done

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

# echo "
# ###################################################
# Making changes to config.txt TODO:
# ###################################################
# "

# echo Making changes to config.txt

# sudo sh -c "echo '\n# EA7KIR Additions' >> /boot/firmware/config.txt"

# echo Disable Wifi
# sudo sh -c "echo 'dtoverlay=disable-wifi' >> /boot/firmware/config.txt"

# echo Disable Bluetooth
# sudo sh -c "echo 'dtoverlay=disable-bt' >> /boot/firmware/config.txt"

echo "
###################################################
Making changes to .profile
###################################################
"

sudo sh -c "echo '\n# EA7KIR Additions' >> /home/pi/.profile"

# echo Disbale Screen Blanking in .profile
# echo -e 'export DISPLAY=:0; xset s noblank; xset s off; xset -dpms' >> /home/pi/.profile

echo Adding go path to .profile
echo -e 'export PATH=$PATH:/usr/local/go/bin' >> /home/pi/.profile

echo "
###################################################
Installing Go $GOVERSION
###################################################
"

GOFILE=go$GOVERSION.linux-arm64.tar.gz
cd /usr/local
sudo wget https://go.dev/dl/$GOFILE
sudo tar -C /usr/local -xzf $GOFILE
cd

echo "
###################################################
Installing gioui dependencies
###################################################
"

sudo apt -y install pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev

echo "
###################################################
Installing gioui tools $GIOUIVERSION
###################################################
"

/usr/local/go/bin/go install gioui.org/cmd/gogio@$GIOUIVERSION

echo "
###################################################
Installing IIO devices
###################################################
"

sudo apt -y install libiio-utils

echo "
###################################################
Installing sshpass
###################################################
"

sudo apt -y install sshpass

echo "
###################################################
Installing plutosdr_scripts/master/ssh_config to /home/pi/.ssh/config
###################################################
"

#mkdir /home/pi/.ssh
wget https://raw.githubusercontent.com/analogdevicesinc/plutosdr_scripts/master/ssh_config -O ~/.ssh/config

# echo "
# ###################################################
# Installing plutosdr-fw/master/scripts/53-adi-plutosdr-usb.rules to /etc/udev/rules.d/
# ###################################################
# "

# sudo wget https://raw.githubusercontent.com/analogdevicesinc/plutosdr-fw/master/scripts/53-adi-plutosdr-usb.rules -O /etc/udev/rules.d/53-adi-plutosdr-usb.rules
# sudo udevadm control --reload-rules

echo "
###################################################
Copying q100transmitter.service
###################################################
"

cd /home/pi/Q100/q100transmitter/etc
sudo cp q100transmitter.service /etc/systemd/system/
sudo chmod 644 /etc/systemd/system/q100transmitter.service
sudo systemctl daemon-reload
cd

echo "
###################################################
Configure routing
###################################################
"

# # Editing /etc/network/interfaces
# TXT="
# auto eth1
#     iface eth1 inet static
#         address 192.168.3.10
#         netmask 255.255.255.0
# "
# echo "$TXT" | sudo tee --append /etc/network/interfaces

# Enable port forwarding
sudo sysctl net.ipv4.ip_forward # check if port forwarding is enabled
sudo sed -i 's/#net.ipv4.ip_forward=1/net.ipv4.ip_forward=1/g' /etc/sysctl.conf # change to 1
sudo sysctl -p # confirm port forwarding is active
lsmod | grep nf # check if nftables is running as a kernel module
sudo systemctl enable nftables
sudo systemctl start nftables

# Adding table
sudo nft list ruleset # check which rules are active
sudo nft flush ruleset # to start over
sudo nft add table nat
sudo nft 'add chain nat postrouting { type nat hook postrouting priority 100 ; }'
sudo sudo nft add rule nat postrouting masquerade
sudo nft 'add chain nat prerouting { type nat hook prerouting priority -100; }'

# Forward the ENCODER streams to the PLUTO
sudo nft add rule nat prerouting iif $ENCODER udp dport 7272 dnat to 192.168.2.1
sudo nft add rule nat prerouting iif $ENCODER udp dport 8282 dnat to 192.168.2.1

# Enable access to ENCODER from the LAN during debug/development
#    ie. access as: http://txtouch.local:8083
sudo nft add rule nat prerouting iif $LAN tcp dport 8083 dnat to 192.168.3.1:80

# Enable access to PLUTO from the LAN during debug/development
#    ie. access as: http://txtouch.local:8082
sudo nft add rule nat prerouting iif $LAN tcp dport 8082 dnat to 192.168.2.1:80

# Checking the rules
sudo nft list ruleset

# Making the rules persist
sudo cp /etc/nftables.conf /etc/nftables.backup
sudo nft list ruleset | sudo tee /etc/nftables.conf

echo "
###################################################
Bring up the encoder using nmcli
###################################################
"

# This appears to work
sudo nmcli con mod Wired\ connection\ 3 ipv4.addresses 192.168.3.10/24
sudo nmcli con mod Wired\ connection\ 3 ipv4.gateway 192.168.3.0
#sudo nmcli con mod Wired\ connection\ 3 ipv4.dns 8.8.8.8
sudo nmcli con mod Wired\ connection\ 3 ipv4.method manual
sudo nmcli con up Wired\ connection\ 3

echo "
###################################################
Prevent this script from being executed again
###################################################
"

chmod -x /home/pi/Q100/q100transmitter/etc/install.sh

echo "
INSTALL HAS COMPLETED

    AFTER REBOOTING:

    Ues your finger to configure some Desktop settings:

    Appearance Settings
	    Disable Wastebasket & External Disks
    Raspberry Pi Configuration
	    System set Network at Boot to ON

    Then login from your PC, Mac, or Linux computer

    ssh pi@txtouch.local

    IMPORTANT: login to the Pluto just once to authenticate
        using passwod 'analog', then 'exit'

    ssh plutosdr

    cd Q100/q100transmitter
    go mod tidy
    go build .
    sudo systemctl enable q100transmitter
    sudo systemctl start q100transmitter

"

while true; do
    read -p "I have read the above, so continue (y/n)? " answer
    case ${answer:0:1} in
        y|Y ) break;;
        n|N ) exit;;
        * ) echo "Please answer yes or no.";;
    esac
done

sudo reboot
