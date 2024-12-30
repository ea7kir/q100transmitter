#!/bin/bash

# take them all  down
sudo nmcli con down "Wired connection 1"
sudo nmcli con down "Wired connection 2"
sudo nmcli con down "Wired connection 3"

# eth0
sudo nmcli con mod "Wired connection 1" ipv4.gateway 192.168.1.1 ipv4.method auto

# Pluto
sudo nmcli con mod "Wired connection 2" ipv4.addresses 192.168.2.10/24 ipv4.method manual ipv6.method ignore
sudo nmcli con mod "Wired connection 2" ipv4.gateway 192.168.2.10

# Encoder
sudo nmcli con mod "Wired connection 3" ipv4.addresses 192.168.3.10/24 ipv4.method manual ipv6.method ignore
sudo nmcli con mod "Wired connection 3" ipv4.gateway 192.168.3.10

# bring them up in order
sudo nmcli con up "Wired connection 1"
sudo nmcli con up "Wired connection 2"
sudo nmcli con up "Wired connection 3"
