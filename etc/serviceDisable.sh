#!/bin/bash

# stop and disable the systemct service

QAPP=q100transmitter

QSERVICE=$QAPP.service

sudo systemctl stop $QSERVICE
sudo systemctl disable $QSERVICE

