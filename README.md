# Q-100 Transmitter
Control and monitor a DATV transmiter with a touch screen.

$${\color{red}WARNING:\space ALL\space DEVELOPMENT\space TAKES\space PLACE\space ON\space THE\space MAIN\space BRANCH}$$

## Hardware
- Raspberry Pi 4B with 4GB RAM (minimum)
- Raspberry Pi Official 7" Touch Screen
- Analog Devices PlutoSDR Revision D
- Analog Devices HMC349 RF Switch 
- Analog Devices EVAL-CN0417-EBZ 2.4GHz RF Power Amplifier
- DDMALL HEV-10 HDMI Video Encoder

**A keyboard and mouse are not required at any time**
## Connections
TODO: add schenatics
## Installing
NOTE: CURRENTLY REQUIRES PI OS BULLSEYE 64-BIT (FULL DESKTOP VERSION)

### Using Raspberry Pi Imager:
```
CHOOSE OS: Raspberry Pi OS (other) -> Raspberry Pi OS (64-bit)

CONFIGURE:
	Set hostname:			txtouch
	Enable SSH
		Use password authentication
	Set username and password
		Username:			pi
		Password: 			<password>
	Set locale settings
		Time zone:			<Europe/Madrid> # or wherever you are
		Keyboard layout:	<us>
	Eject media when finished
SAVE and WRITE
```

Insert the card into the Raspberry Pi and switch on

WARNING: the Pi may reboot during the install, so please allow it to complete

### Remote login from a Mac, PC or Linux host
```
ssh pi@txtouch.local

mkdir Q100
cd Q100
git clone https://github.com/ea7kir/q100transmitter.git

cd q100transmitter/etc
chmod +x install.sh
./install.sh
```

THEN FOLLOW THE INSTRUCTIONS TO CONFIGURE THE DESKTOP

## License
Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see https://www.gnu.org/licenses/.
