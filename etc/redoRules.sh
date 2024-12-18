#!/bin/bash

################ NEW VERSION ###########################################

# echo "
# ###################################################
# NetworkManager configure for TxTouch
# 	will start in 5 seconds
# ###################################################
# "

# sleep 5

# echo "
# ###################################################
# Get the Encoder connected
#     assuming it's connected to eth2
# ###################################################
# "
# nmcli dev
# sudo nmcli con add con-name Wired\ connection\ 3 type ethernet ifname eth2 ipv4.method manual ipv4.address 192.168.3.10/24
# nmcli dev

# echo "
# ###################################################
# Rename the known connection names
# 	and persist on boot
# ###################################################
# "

# nmcli dev
# sudo nmcli con mod Wired\ connection\ 1 connection.id Lan
# sudo nmcli con mod Wired\ connection\ 2 connection.id Pluto
# sudo nmcli con mod Wired\ connection\ 3 connection.id Encoder
# nmcli dev

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

echo "
###################################################
Add my-nat network translation table
###################################################
"

sudo nft list ruleset # check which rules are active
sudo nft flush ruleset # to start over
sudo nft add table my-nat
sudo nft 'add chain my-nat postrouting { type nat hook postrouting priority 100 ; }'
sudo sudo nft add rule my-nat postrouting masquerade
sudo nft 'add chain my-nat prerouting { type nat hook prerouting priority -100; }'

sudo nft list ruleset # check which rules are active

echo "
###################################################
Enable access to PLUTO from the LAN during debug/development
	ie. access as: http://txtouch.local:8082
###################################################
"

sudo nft add rule my-nat prerouting iif eth0 tcp dport 8082 dnat to 192.168.2.1:80

echo "
###################################################
Enable access to ENCODER from the LAN during debug/development
	ie. access as: http://txtouch.local:8083
###################################################
"

sudo nft add rule my-nat prerouting iif eth0 tcp dport 8083 dnat to 192.168.3.1:80

echo "
###################################################
Forward the ENCODER streams to the PLUTO
###################################################
"

sudo nft add rule my-nat prerouting iif eth2 udp dport 7272 dnat to 192.168.2.1
sudo nft add rule my-nat prerouting iif eth2 udp dport 8282 dnat to 192.168.2.1

echo "
###################################################
Making the rules persist
    do not save the backup to /etc/nftables.backup
###################################################
"

# sudo cp /etc/nftables.conf /home/pi/nftables.backup # save a backup copy
sudo cp /etc/nftables.conf /home/pi/nftables.backup-2 # save a backup copy
sudo nft list ruleset | sudo tee /etc/nftables.conf # persist the table

sleep 5

#############END OF NEW #################################################