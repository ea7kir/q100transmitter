# Configuring the Network

An attempt to use the lastest network tools, such as ip and nftables.

NOTE: With the PLUTO connected to the TOP USB socket.

NOTE: With the ENCODER connected via a REALTEK OTG to the BOTTOM USB socket.

## Before changing anything

pi@txtouch:~ $ `lsusb`

    Bus 002 Device 002: ID 0bda:8153 Realtek Semiconductor Corp. RTL8153 Gigabit Ethernet Adapter
    Bus 002 Device 001: ID 1d6b:0003 Linux Foundation 3.0 root hub
    Bus 001 Device 003: ID 0456:b673 Analog Devices, Inc. PlutoSDR (ADALM-PLUTO)
    Bus 001 Device 002: ID 2109:3431 VIA Labs, Inc. Hub
    Bus 001 Device 001: ID 1d6b:0002 Linux Foundation 2.0 root hub

Therefore, Bus 002 Device 002 is the ENCODER, and Bus 001 Device 003 is the PLUTO.

pi@txtouch:~ $ `ip addr`

    1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
        link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
        inet 127.0.0.1/8 scope host lo
           valid_lft forever preferred_lft forever
        inet6 ::1/128 scope host 
           valid_lft forever preferred_lft forever
    2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
        link/ether e4:5f:01:9f:14:c3 brd ff:ff:ff:ff:ff:ff
        inet 192.168.1.37/24 brd 192.168.1.255 scope global dynamic noprefixroute eth0
           valid_lft 43011sec preferred_lft 37611sec
        inet6 fe80::20f5:e7a1:7a99:852/64 scope link 
           valid_lft forever preferred_lft forever
    3: eth1: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
        link/ether 00:e0:4c:88:44:bc brd ff:ff:ff:ff:ff:ff
        inet 169.254.142.42/16 brd 169.254.255.255 scope global noprefixroute eth1
           valid_lft forever preferred_lft forever
        inet6 fe80::a5ee:8e71:1497:93fb/64 scope link 
           valid_lft forever preferred_lft forever
    4: wlan0: <BROADCAST,MULTICAST> mtu 1500 qdisc noop state DOWN group default qlen 1000
        link/ether e4:5f:01:9f:14:c4 brd ff:ff:ff:ff:ff:ff
    5: eth2: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UNKNOWN group default qlen 1000
        link/ether 00:e0:22:6f:28:9a brd ff:ff:ff:ff:ff:ff
        inet 192.168.2.10/24 brd 192.168.2.255 scope global dynamic noprefixroute eth2
           valid_lft 863807sec preferred_lft 755807sec
        inet6 fe80::c7f3:72d4:8a12:1c23/64 scope link 
           valid_lft forever preferred_lft forever

Thus showing:

    eth0 is the LAN with ip 102.168.1.38/24 (assigned by DHCP from my router)
    eth1 is the ENCODER with ip 169.254.142.42/16 (a Link-local address assigned by linux)
    eth2 is the PLUTO with ip 192.168.2.10/24 (assigned by the pluto)

and:

pi@txtouch:~ $ `ip route`

    default via 192.168.1.1 dev eth0 proto dhcp src 192.168.1.37 metric 202 
    169.254.0.0/16 dev eth1 scope link src 169.254.142.42 metric 203 
    192.168.1.0/24 dev eth0 proto dhcp scope link src 192.168.1.37 metric 202 
    192.168.2.0/24 dev eth2 proto dhcp scope link src 192.168.2.10 metric 205 

## Configure the ENCODER

As the ENCODER is unable to run as a DHCP client, I decided to mimick the way Pluto behaves...

I will configure the ENCODER to be 192.1.3.1 and I'll assign 192.168.3.10 to eth1.

Edit /etc/network/interfaces

`sudo nano /etc/network/interfaces`

and add this:

    auto eth1
        iface eth1 inet static
	        address 192.168.3.10
	        netmask 255.255.255.0

## Prepare IP Routing

It may be neccessary to enable port forwarding if it is not already running.

See: `https://wiki.nftables.org/wiki-nftables/index.php/Performing_Network_Address_Translation_(NAT)`

### Enable port forwarding:

`sudo sysctl net.ipv4.ip_forward` # check if port forwarding is enabled

`sudo sed -i 's/#net.ipv4.ip_forward=1/net.ipv4.ip_forward=1/g' /etc/sysctl.conf` # change to 1

`sudo sysctl -p` # confirm port forwarding is active

`lsmod | grep nf` # check if nftables is running as a kernel module

### If it is not running:

`sudo systemctl enable nftables`

`sudo systemctl start nftables`

## Continue:

`sudo nft list ruleset` # check which rules are active

`sudo nft flush ruleset` # to start over

`sudo nft add table nat`

`sudo nft 'add chain nat postrouting { type nat hook postrouting priority 100 ; }'`

`sudo sudo nft add rule nat postrouting masquerade`

`sudo nft 'add chain nat prerouting { type nat hook prerouting priority -100; }'`

## Forward the ENCODER streams to the PLUTO

`sudo nft add rule nat prerouting iif eth1 udp dport 7272 dnat to 192.168.2.1`

`sudo nft add rule nat prerouting iif eth1 udp dport 8282 dnat to 192.168.2.1`

## Access the ENCODER from the LAN during debug/development 

`sudo nft add rule nat prerouting iif eth0 tcp dport 8083 dnat to 192.168.3.1:80`

I.e. access as: `http://txtouch.local:8083` 

## Access the PLUTO from the LAN during debug/development 

`sudo nft add rule nat prerouting iif eth0 tcp dport 8082 dnat to 192.168.2.1:80`

I.e. access as: `http://txtouch.local:8082`

## Check the rules

`sudo nft list ruleset`

## Make the rules persist

`sudo cp /etc/nftables.conf /etc/nftables.backup`

`sudo  nft list ruleset | sudo tee /etc/nftables.conf`
