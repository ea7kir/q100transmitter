# IP Routing

DDMALL HEV-10 Encoder IP Address is set to 192.168.3.1 username/password admin/admin, and sends to 192.168.3.10

PLUTO IP Address is set to 192.168.2.1 username/password root/analog, and sends to 192.168.2.10

Some usefull USB commands:

    lsusb
    usb-devices
    dmesg | grep usb

Current IP/Routing giving web access to the ENCODER and PLUTO over the LAN

## Assign a new IP to eth1

```
sudo nano /etc/network/interfaces
```

and add this:

```
auto eth1
iface eth1 inet static
	address 192.168.3.10
	netmask 255.255.255.0
	gateway 192.168.3.0
	dns-nameservers 8.8.8.8 8.8.4.4
```

## Port Forwarding

