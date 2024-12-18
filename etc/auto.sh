#!/bin/bash


# lan
sudo nmcli con down Wired\ connection\ 1
sudo nmcli con mod Wired\ connection\ 1 connectio.id myLAN
sudo nmcli con up myLAN

# pluto
sudo nmcli con down Wired\ connection\ 2
sudo nmcli con mod Wired\ connection\ 2 connectio.id myPluto
sudo nmcli con up myPluto

#encoder
sudo nmcli con down Wired\ connection\ 3
sudo nmcli con mod Wired\ connection\ 3 connectio.id myEncoder
sudo nmcli con up myEncoder


# pluto
sudo nmcli con down Wired\ connection\ 2
sleep 5 # allow the router dns service to catch up
sudo nmcli con mod Wired\ connection\ 2 ipv4.addresses 192.168.2.10/24
sleep 5 # allow the router dns service to catch up
#sudo nmcli con mod Wired\ connection\ 2 ipv4.gateway 192.168.2.0
sudo nmcli con mod Wired\ connection\ 2 ipv4.method manual
sleep 5 # allow the router dns service to catch up
sudo nmcli con up Wired\ connection\ 2
sleep 5 # allow the router dns service to catch up

# encoder
sudo nmcli con down Wired\ connection\ 3
sleep 5 # allow the router dns service to catch up
sudo nmcli con mod Wired\ connection\ 3 ipv4.addresses 192.168.3.10/24
sleep 5 # allow the router dns service to catch up
#sudo nmcli con mod Wired\ connection\ 3 ipv4.gateway 192.168.3.0
sudo nmcli con mod Wired\ connection\ 3 ipv4.method manual
sleep 5 # allow the router dns service to catch up
sudo nmcli con up Wired\ connection\ 3
sleep 5 # allow the router dns service to catch up

sudo nmcli con mod Wired\ connection\ 1 connectio.id myLAN
sudo nmcli con mod Wired\ connection\ 2 connectio.id myPluto
sudo nmcli con mod Wired\ connection\ 3 connectio.id myEncoder