# HEV-10 Current Settings

iMac has 2nd IP 192.168.3.3

HEV-10 IP address is 192.168.3.1 and sending UDP stream to 192.168.3.10

NOTE:ignore items in (bracket) as they are not selectable

## Status Tab
### Input
```
Video Input: 1920x1080p@25Hz
```
### Primary Stream
```
       Audio Encoding Type: AAC
        Audio Bitrate(bps): 64000
       Video Encoding Type: H.265
       Video Encoding Size: 1280*720
       Video Bitrate(Kbps): 350
                 RTSP URL1: Disabled
             RTSP URL2(TS): Disabled
UDP(unicast/multicast) URL: udp://192.168.3.10:8282
                       SRT: Disabled
```
### Secondary Stream
```
       Audio Encoding Type: AAC
        Audio Bitrate(bps): 64000
       Video Encoding Type: H.265
       Video Encoding Size: 1280*720
       Video Bitrate(Kbps): 330
                 RTSP URL1: Disabled
             RTSP URL2(TS): Disabled
UDP(unicast/multicast) URL: Disabled
                       SRT: Disabled
```
## Audio Tab
```
Encoding Type: ACC
Bit Rate(bps): 64000
```
## Video Tab
### Primary Stream
#### Encoding setting
```
           Input Status: 1920x1080p@25Hz
          Encoding Type: H.265
        Frame Rate(fps): 25
Key Frame Interval(GOP): 50
               Bit Rate: Manual Input   350
      Output Resolution: 1250*720
    Compression Profile: (Main)
        Bitrate Control: CBR
```
### Secondary Stream
#### Encoding setting
```
               Encoding: OFF
           Input Status: 1920x1080p@25Hz
          Encoding Type: (H.265)
        Frame Rate(fps): (25)
Key Frame Interval(GOP): 50
               Bit Rate: (Manual Input)   330
      Output Resolution: (1250*720)
    Compression Profile: (Main)
        Bitrate Control: (CBR)
```
## Network Tab
### Local
```
    IP Mode: Static
 IP Address: 192.168.3.1
Subnet Mask: 255.255.255.0
    Gateway: 192.168.1.1
        DNS: 8.8.8.8
```
### UPnP
```
UPnP: OFF
```
## DDNS
```
DDNS: OFF
```
## RTMP Tab
### Primary Stream
#### No.1
```
   Server URL: //undefined
   Stream Key:
       Status: disconnected
Athentication:
   User Name:
    Password:
```
### Secondary Stream
#### No.1
```
   Server URL:
   Stream Key:
       Status: disconnected
Athentication:
   User Name:
    Password:
```
## Misc Stream Tab
### Primary Stream
```
Stream Protocol: UDP

       Protocol: TS over UDP (Unicast/Multicast)
 Destination IP: 192.168.3.10
           Port: 8282
            TTL: 64
```
### Secondary Stream
```
* Please turn on video encoding for secondary stream firstly.
Stream Protocol: (RTSP)

     URL.1 Port: (8553)
    URL1: (rtsp://192.168.3.1:8553/live1)
URL.2 (TS) Port: (553)
      URL2 (TS): rtsp://192.168.3.1:553/live1
      Multicast: (Disabled)
   Multicast IP: (239.255.17.18)
 Multicast Port: (8530)
            TTL: (64)
 Authentication: (Disabled)
``
## System Tab
### Device Info
```
       Serial No.: ENF2202Y30W29010
      Device Name: HEV-10
Firmeware Version: HEV10-V1.1.2
 Firmeare Uograde: Choose File
```
### User Configure
```
       Old User Name:
        Old Password:
       New User Name:
        New Password:
Confirm New Password:
```