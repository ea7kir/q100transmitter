#!/bin/bash

##  NEDS UPDATING

# Install Q100 Transmitter on txtouch.local
# Orignal design by Michael, EA7KIR

# CONFIFIGURATION
GOVERSION=1.23.4
GIOUIVERSION=v0.7.1

# This is what we hope for if all goes well

# nmcli device
# DEVICE         TYPE      STATE                   CONNECTION         
# eth0           ethernet  connected               Wired connection 1 
# eth1           ethernet  connected               Wired connection 2 
# eth2           ethernet  connected               Wired connection 3 
# lo             loopback  connected (externally)  lo                 
# wlan0          wifi      unavailable             --                 

whoami | grep -q pi
if [ $? != 0 ]; then
  echo Install must be performed as user pi
  exit
fi

hostname | grep -q TxTouch
if [ $? != 0 ]; then
  echo Install must be performed on host TxTouch
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

# echo "
# ###################################################
# Update Pi OS
# ###################################################
# "

sudo apt update
sudo apt -y full-upgrade
sudo apt -y autoremove
sudo apt clean

# echo "
# ###################################################
# Making changes to config.txt TODO:
# ###################################################
# "

# TODO: find the coorect way of doing these

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

# TODO: find the coorect way of doing this
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

sudo apt -y install gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev

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
    first time use: login as: ssh root@192.168.2.1
    later logins can use: ssh plutosdr or: sshpass -p analog ssh plutosdr
###################################################
"

# see: https://wiki.analog.com/university/tools/pluto/drivers/linux
# TODO: copy from the etc folder
mkdir /home/pi/.ssh
touch /home/pi/.ssh/known_hosts
wget https://raw.githubusercontent.com/analogdevicesinc/plutosdr_scripts/master/ssh_config -O /home/pi/.ssh/config

# NOTE: why is this still here?
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
Re-Configure eth1 and eth2
    leaving eth0 as default (using DHCP)
###################################################
"

# turn of bluetooth and wifi
sudo nmcli radio wifi off
sudo rfkill block all

# take them all down, so we can them bring them
# back up in order of priority
sudo nmcli con down "Wired connection 1"
sudo nmcli con down "Wired connection 2"
sudo nmcli con down "Wired connection 3"

# eth0
# sudo nmcli con mod "Wired connection 1" ipv4.gateway 192.168.1.1 ipv4.method auto # fails most times
sudo nmcli con mod "Wired connection 1" connection.autoconnect-priority 903

# Pluto
sudo nmcli con mod "Wired connection 2" ipv4.addresses 192.168.2.10/24 ipv4.method manual ipv6.method ignore
sudo nmcli con mod "Wired connection 2" ipv4.gateway 192.168.2.10
sudo nmcli con mod "Wired connection 2" connection.autoconnect-priority 902

# Encoder
sudo nmcli con mod "Wired connection 3" ipv4.addresses 192.168.3.10/24 ipv4.method manual ipv6.method ignore
sudo nmcli con mod "Wired connection 3" ipv4.gateway 192.168.3.10
sudo nmcli con mod "Wired connection 3" connection.autoconnect-priority 901

# bring them up in order of priority
sudo nmcli con up "Wired connection 1"
sudo nmcli con up "Wired connection 2"
sudo nmcli con up "Wired connection 3"

sleep 5
sudo nmcli con reload

sleep 5
nmcli dev

echo "
###################################################
Enable port forwarding
###################################################
"

sudo sysctl net.ipv4.ip_forward # check if port forwarding is enabled
sudo sed -i 's/#net.ipv4.ip_forward=1/net.ipv4.ip_forward=1/g' /etc/sysctl.conf # change to 1
sudo sysctl -p # confirm port forwarding is active
lsmod | grep nf # check if nftables is running as a kernel module

echo "
###################################################
Enable nftables
###################################################
"

sudo systemctl enable nftables
sudo systemctl start nftables

sudo systemctl status nftables
lsmod | grep nf # check if nftables is running as a kernel module

echo "
###################################################
Configure routes
###################################################
"

# start over
sudo nft flush ruleset

# Add nat-table network translation table
# NOTE: I don't remeber where I fould this gem
sudo nft add table nat-table
sudo nft 'add chain nat-table postrouting { type nat hook postrouting priority 100 ; }'
sudo nft add rule nat-table postrouting masquerade
sudo nft 'add chain nat-table prerouting { type nat hook prerouting priority -100; }'

# Forward the ENCODER stream to the PLUTO (7272 for H224, 8282 for H265)
sudo nft add rule nat-table prerouting iif eth2 udp dport 7272 dnat to 192.168.2.1
sudo nft add rule nat-table prerouting iif eth2 udp dport 8282 dnat to 192.168.2.1

# Enable HTTP access to PLUTO from the LAN
#   ie. access as: http://txtouch.local:8082
sudo nft add rule nat-table prerouting iif eth0 tcp dport 8082 dnat to 192.168.2.1:80

# Enable HTTP access to ENCODER from the LAN
# 	ie. access as: http://txtouch.local:8083
sudo nft add rule nat-table prerouting iif eth0 tcp dport 8083 dnat to 192.168.3.1:80

sudo nft list ruleset

echo "
###################################################
Persist after booting
###################################################
"

# NOTE: I don't remeber where I fould this gem
sudo cp /etc/nftables.conf /home/pi/nftables.backup # save a backup copy
sudo nft list ruleset | sudo tee /etc/nftables.conf # persist the table

sleep 5

echo "
###################################################
Prevent this script from being executed again
###################################################
"

chmod -x /home/pi/Q100/q100transmitter/etc/install.sh

echo "
###################################################
INSTALL HAS COMPLETED
###################################################
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
