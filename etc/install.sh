#!/bin/bash

# Install Q100 Transmitter on Raspberry Pi 4
# Orignal design by Michael, EA7KIR

GOVERSION=1.21.4

echo WARNING: THIS INSTALL SCRIPT HAS NOT BEEN TESTED

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

while true; do
    read -p "Install q100transmitter using Go version $GOVERSION (y/n)? " answer
    case ${answer:0:1} in
        y|Y ) break;;
        n|N ) exit;;
        * ) echo "Please answer yes or no.";;
    esac
done

mkdir /home/pi/Q100

read -p "Enter your callsign in uppercase and press enter " callsign
echo $callsign > /home/pi/Q100/callsign

echo "\n###################################################\n"

echo Updateing Pi OS
sudo apt update
sudo apt -y full-upgrade
sudo apt -y autoremove
sudo apt clean

echo Running rfkill # not sure if this dupicates config.txt
rfkill block 0
rfkill block 1

echo Making changes to config.txt

echo Disbaling Wifi
echo -e "\ndtoverlay=disable-wifi" >> /boot/config.txt

echo Disbaling Bluetooth
echo -e "\ndtoverlay=disable-bt" >> /boot/config.txt

echo Installing IIO devices
sudo apt install libiio-utils

echo Installing sshpass
sudo apt install sshpass

echo Installing plutosdr_scripts/master/ssh_config to /home/pi/.ssh/config
wget https://raw.githubusercontent.com/analogdevicesinc/plutosdr_scripts/master/ssh_config -O ~/.ssh/config

echo Installing plutosdr-fw/master/scripts/53-adi-plutosdr-usb.rules to /etc/udev/rules.d/
sudo wget https://raw.githubusercontent.com/analogdevicesinc/plutosdr-fw/master/scripts/53-adi-plutosdr-usb.rules -O /etc/udev/rules.d/
sudo udevadm control --reload-rules

echo "\n###################################################\n"

echo Adding go path to .profile
echo -e '\n\nexport PATH=$PATH:/usr/local/go/bin\n\n' >> /home/pi/.profile

echo Installing Go $GOVERSION
GOFILE=go$GOVERSION.linux-arm64.tar.gz
cd /usr/local
sudo wget https://go.dev/dl/$GOFILE
sudo tar -C /usr/local -xzf $GOFILE
cd

echo Installing gioui dependencies
sudo apt install gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev

echo Installing gioui tools
go install gioui.org/cmd/gogio@latest

echo "\n###################################################\n"

echo Copying q100transmitter.service
cd /home/pi/Q100/q100transmitter/etc
sudo cp q100transmitter.service /etc/systemd/system/
sudo chmod 644 /etc/systemd/system/q100transmitter.service
sudo systemctl daemon-reload
cd

###################### begin configure routing ######################

echo Configure routing

echo Editing /etc/network/interfaces
TXT="
auto eth1
    iface eth1 inet static
        address 192.168.3.10
        netmask 255.255.255.0
"
echo "$TXT" | sudo tee --append /etc/network/interfaces

echo Enable port forwarding
sudo sysctl net.ipv4.ip_forward # check if port forwarding is enabled
sudo sed -i 's/#net.ipv4.ip_forward=1/net.ipv4.ip_forward=1/g' /etc/sysctl.conf # change to 1
sudo sysctl -p # confirm port forwarding is active
lsmod | grep nf # check if nftables is running as a kernel module
sudo systemctl enable nftables
sudo systemctl start nftables

echo Adding table
sudo nft list ruleset # check which rules are active
sudo nft flush ruleset # to start over
sudo nft add table nat
sudo nft 'add chain nat postrouting { type nat hook postrouting priority 100 ; }'
sudo sudo nft add rule nat postrouting masquerade
sudo nft 'add chain nat prerouting { type nat hook prerouting priority -100; }'

echo Forward the ENCODER streams to the PLUTO
sudo nft add rule nat prerouting iif eth1 udp dport 7272 dnat to 192.168.2.1
sudo nft add rule nat prerouting iif eth1 udp dport 8282 dnat to 192.168.2.1

echo Enable access to ENCODER from the LAN during debug/development
echo    ie. access as: http://txtouch.local:8083
sudo nft add rule nat prerouting iif eth0 tcp dport 8083 dnat to 192.168.3.1:80

echo Enable access to PLUTO from the LAN during debug/development
echo    ie. access as: http://txtouch.local:8082
sudo nft add rule nat prerouting iif eth0 tcp dport 8082 dnat to 192.168.2.1:80

echo Checking the rules
sudo nft list ruleset

echo Making the rules persist
sudo cp /etc/nftables.conf /etc/nftables.backup
sudo nft list ruleset | sudo tee /etc/nftables.conf

###################### endof configure routing ######################

echo "\n
INSTALL HAS COMPLETED
   after rebooting, build and auto exec...

   It will be neccessary to log into the Pluto once to get things working...

   ssh TODO:

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
