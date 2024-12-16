#!/bin/bash

GOVERSION=1.23.4

GOFILE=go$GOVERSION.linux-arm64.tar.gz

cd /usr/local
sudo wget https://go.dev/dl/$GOFILE
sudo rm -rf /usr/local/go 
sudo tar -C /usr/local -xzf $GOFILE
cd $OLDPWD
go version

