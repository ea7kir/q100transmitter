# ROB-NET-CONFIG
```
From: Michael Naylor <ea7kir@icloud.com>
Subject: Re: NetworkManager
Date: 31 December 2024 at 17:19:15 CET
To: Rob Vasak <mailrob@telstra.com>

Hi Rob,

Eureka it’s working!

You have exceeded the call to duty.  I can’t thank you enough. This is World-Wide Amateur DATV at it’s best.

This time I managed to document the steps taken (below).  The text mentions file "fix-routes.sh" which I include a copy for reference.

I look forward to meeting you on QO-100 one day.  Some over here are experimenting with SRT.  I intend to build SRT into the fake repeater I’m also developing.  I can get back to that now.

Happy New Year,
Michael Naylor - EA7KIR
```
### include: fix-routes.sh


Following Rob’s December 31st email advice - 31 December

## 1. Check eth2 Initialization

After reboot, check the status of eth2:

nmcli dev
```
pi@TxTouch:~ $ nmcli dev
DEVICE  TYPE      STATE                                  CONNECTION         
eth0    ethernet  connected                              Wired connection 1 
eth1    ethernet  connected                              Wired connection 2 
lo      loopback  connected (externally)                 lo                 
eth2    ethernet  connecting (getting IP configuration)  Wired connection 3 
wlan0   wifi      unavailable
```
ip addr show eth2
```
pi@TxTouch:~ $ ip addr show eth2
4: eth2: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 00:0e:c6:75:b5:c7 brd ff:ff:ff:ff:ff:ff
    inet6 fe80::da89:9e58:788f:a790/64 scope link noprefixroute 
       valid_lft forever preferred_lft forever
```

If eth2 is missing or not configured with a link-local address (169.254.x.x), it means the interface isn’t being properly initialized.

## 2. Force Initialization of eth2

Ensure that eth2 is assigned an IP address, even if the Encoder isn't powered on.

Option A: Assign a Static IP

You can configure a static IP address for eth2 in /etc/network/interfaces or using nmcli:

sudo nmcli con mod "Wired connection 3" ipv4.addresses 192.168.3.10/24 ipv4.method manual
sudo nmcli con up "Wired connection 3"

```
pi@TxTouch:~ $ sudo nmcli con mod "Wired connection 3" ipv4.addresses 192.168.3.10/24 ipv4.method manual
pi@TxTouch:~ $ sudo nmcli con up "Wired connection 3"
Connection successfully activated (D-Bus active path: /org/freedesktop/NetworkManager/ActiveConnection/9)

pi@TxTouch:~ $ sudo reboot

pi@TxTouch:~ $ nmcli dev
DEVICE  TYPE      STATE                   CONNECTION         
eth0    ethernet  connected               Wired connection 1 
eth1    ethernet  connected               Wired connection 2 
eth2    ethernet  connected               Wired connection 3 
lo      loopback  connected (externally)  lo                 
wlan0   wifi      unavailable             --                 

pi@TxTouch:~ $ ip addr show eth2
4: eth2: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 00:0e:c6:75:b5:c7 brd ff:ff:ff:ff:ff:ff
    inet 192.168.3.10/24 brd 192.168.3.255 scope global noprefixroute eth2
       valid_lft forever preferred_lft forever
    inet6 fe80::da89:9e58:788f:a790/64 scope link noprefixroute 
       valid_lft forever preferred_lft forever
```

## So far so good. On to part 3. Update nftables Rules 

This is a fresh OS install with only the sefault /etc/nftables.conf….

```
pi@TxTouch:~ $ cat /etc/nftables.conf
#!/usr/sbin/nft -f

flush ruleset

table inet filter {
        chain input {
                type filter hook input priority filter;
        }
        chain forward {
                type filter hook forward priority filter;
        }
        chain output {
                type filter hook output priority filter;
        }
}
```

So I need to configure my own routing using my “fix-routes.sh” file…

```
pi@TxTouch:~ $ cat /etc/nftables.conf
table ip nat-table {
	chain postrouting {
		type nat hook postrouting priority srcnat; policy accept;
		masquerade
	}

	chain prerouting {
		type nat hook prerouting priority dstnat; policy accept;
		iif "eth2" udp dport 7272 dnat to 192.168.2.1
		iif "eth2" udp dport 8282 dnat to 192.168.2.1
		iif "eth0" tcp dport 8082 dnat to 192.168.2.1:80
		iif "eth0" tcp dport 8083 dnat to 192.168.3.1:80
	}
}

pi@TxTouch:~ $ sudo sed -i 's/iif "eth2"/iifname "eth2"/g' /etc/nftables.conf
pi@TxTouch:~ $ cat /etc/nftables.conf 
table ip nat-table {
	chain postrouting {
		type nat hook postrouting priority srcnat; policy accept;
		masquerade
	}

	chain prerouting {
		type nat hook prerouting priority dstnat; policy accept;
		iifname "eth2" udp dport 7272 dnat to 192.168.2.1
		iifname "eth2" udp dport 8282 dnat to 192.168.2.1
		iif "eth0" tcp dport 8082 dnat to 192.168.2.1:80
		iif "eth0" tcp dport 8083 dnat to 192.168.3.1:80
	}
}
```

## So far so good. On to part 4. Persist Interface Configurartion
```
pi@TxTouch:~ $ sudo nmcli con up "Wired connection 3"
Connection successfully activated (D-Bus active path: /org/freedesktop/NetworkManager/ActiveConnection/5)
```
## So far so good. On to part 5. Restart the System and Verify

```
pi@TxTouch:~ $ sudo reboot

pi@TxTouch:~ $ nmcli dev
DEVICE  TYPE      STATE                   CONNECTION         
eth0    ethernet  connected               Wired connection 1 
eth1    ethernet  connected               Wired connection 2 
eth2    ethernet  connected               Wired connection 3 
lo      loopback  connected (externally)  lo                 
wlan0   wifi      unavailable             --                 

pi@TxTouch:~ $ ip addr show eth2
4: eth2: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 00:0e:c6:75:b5:c7 brd ff:ff:ff:ff:ff:ff
    inet 192.168.3.10/24 brd 192.168.3.255 scope global noprefixroute eth2
       valid_lft forever preferred_lft forever
    inet6 fe80::da89:9e58:788f:a790/64 scope link noprefixroute 
       valid_lft forever preferred_lft forever

pi@TxTouch:~ $ sudo nft list ruleset
table ip nat-table {
	chain postrouting {
		type nat hook postrouting priority srcnat; policy accept;
		masquerade
	}

	chain prerouting {
		type nat hook prerouting priority dstnat; policy accept;
		iifname "eth2" udp dport 7272 dnat to 192.168.2.1
		iifname "eth2" udp dport 8282 dnat to 192.168.2.1
		iif "eth0" tcp dport 8082 dnat to 192.168.2.1:80
		iif "eth0" tcp dport 8083 dnat to 192.168.3.1:80
	}
}
```

## So far so good. Now to Test the stream from the Encoder to Pluto

1. I can access the Encoder and Pluto from Safari on my iMac
2. The Pluto Transport stream anaysis page show activity, but there is no video
3. I can ping google.com from the Pi (on previous attemps this hasn’t always been possible)
4. To test some more, I’ll need to install my q100transmitter application which requires an OS upgrade
5. Afer upgrading I ran the verififacation steps again and the result were the same.
6. installed q100transmitter and all dependences, compiled and ran - IT WORKS !
7. Lucky Seven, but will it survive a shutdown and run from my systemctl service? YES !

end.
