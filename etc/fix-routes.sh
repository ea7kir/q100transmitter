#!/bin/bash

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
