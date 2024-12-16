#!/bin/bash

# install, enable and start the systemct service

QAPP=q100transmitter

QPATH=/home/pi/Q100/$QAPP/etc

QSERVICE=$QAPP.service

sudo cp $QPATH/$QSERVICE /etc/systemd/system/
sudo chmod 644 /etc/systemd/system/$QSERVICE
sudo systemctl daemon-reload

sudo systemctl enable $QSERVICE
sudo systemctl start $QSERVICE

