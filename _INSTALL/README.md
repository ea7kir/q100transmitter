## Hardware Connections

**IMPORTANT:** The following peripherals MUST be configured, connected and powered before installing.

- Raspberry Pi 4B with 4GB RAM (minimum)
- Raspberry Pi Official 7" Touch Screen
- Analog Devices PlutoSDR Revision D
- Analog Devices HMC349 RF Switch 
- Analog Devices EVAL-CN0417-EBZ 2.4GHz RF Power Amplifier
- DDMALL HEV-10 HDMI Video Encoder

In addtion, a hefty 5v PSU will be required, or, for example...

- 12v PSU
- 5v Buck Converter

## Installing Pi OS

NOTE: CURRENTLY REQUIRES PI OS BULLSEYE 64-BIT (FULL DESKTOP VERSION)

### Using Raspberry Pi Imager:

```
CHOOSE OS: Raspberry Pi OS (other) -> Raspberry Pi OS (64-bit)

CONFIGURE:
	Set hostname:			q100transmitter
	Enable SSH
		Use password authentication
	Set username and password
		Username:			pi
		Password: 			<password>
	Set locale settings
		Time zone:			<Europe/Madrid>
		Keyboard layout:	<us>
	Eject media when finished
SAVE and WRITE
```

Insert the card into the Raspberry Pi and switch on

WARNING: the Pi may reboot during the install, so please allow it to complete

## Remote login from a Mac, PC or Linux host

```
ssh pi@q100transmitter.local
```

Enable GPIO and Expand the filesystem - using raspi-config

```
sudo raspi-config
```

It's also possible to remove packages that won't be needed - using any of these commands

```
sudo apt purge -y [packagename]
sudo apt autoremove
sudo apt autoclean
```

The difference between “autoremove” and “autoclean” is that “autoremove” removes all unused packages from the system, while “autoclean” removes unused packages from the source repository.

## Clone the repro

```
mkdir Q100
cd Q100
git clone https://github.com/ea7kir/q100transmitter.git
```

Execute the install scripts in order

```
cd
cp /home/pi/Q100/q100/q100transmitter/_INSTALL/install_* .
chmod +x install_*
./install_1.sh
```

This will fire off a sequence of events invoving more than one reboot, so please allow all scripts to complete

## What next?

Remove the install scripts from the home directory

```
cd
rm install_*
```

It will be neccessary to login in to the Pluto at least once (password is `analog`) in order to establish a local certificate, otherwise Tune will not work correctly.
```
ssh root@pluto.local
```