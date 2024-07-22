#!/bin/bash

# Update Q100 Transmitter on Raspberry Pi 4
# Orignal design by Michael, EA7KIR

GOVERSION=1.22.5

whoami | grep -q pi
if [ $? != 0 ]; then
  echo Update must be performed as user pi
  exit
fi

hostname | grep -q txtouch
if [ $? != 0 ]; then
  echo Update must be performed on host rxtouch
  exit
fi

while true; do
   read -p "Update q100transmitter using Go version $GOVERSION (y/n)? " answer
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

sudo apt update
sudo apt -y full-upgrade
sudo apt -y autoremove
sudo apt clean

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

# sudo apt install gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev
sudo apt -y install pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev

echo "###################################################
Installing gioui tools
###################################################
"

/usr/local/go/bin/go install gioui.org/cmd/gogio@latest

echo "
###################################################
Copying q100transmitter.service
###################################################
"

cd /home/pi/Q100/q100transmitter/etc
sudo cp q100receiver.service /etc/systemd/system/
sudo chmod 644 /etc/systemd/system/q100transmiter.service
sudo systemctl daemon-reload
cd

echo "
UPDATE HAS COMPLETED

    AFTER REBOOTING...

    Login from your PC, Mc, or Linux computer

    ssh pi@txtouch.local

    and either execute the following commands
    
    cd Q100/q100transmitter
    go mod tidy
    go build .
    sudo systemctl enable q100transmitter
    sudo systemctl start q100transmitter

    or just

    cd Q100/q100transmitter
    go mod tidy
    go run .

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
