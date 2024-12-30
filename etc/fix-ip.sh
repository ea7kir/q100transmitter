#!/bin/bash

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

