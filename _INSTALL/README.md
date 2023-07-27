## Hardware Connections

TODO: list pin connections and refer to drawings

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

Clone the repro

```
mkdir Q100
cd Q100
git clone https://github.com/ea7kir/q100transmitter.git
```

Execute the install scripts

```
cd
cp /home/pi/_INSTALL/install_* .
chmod +x install_*
./install_1.sh
```

This will fire off a sequnce of events incloving more than one reboot, so please allow it to complete

## More to follow goes here

NOTE: it will be neccessary to login in to the Pluto at least once (password is `analog`) in order to establish a local certificate, otherwise Tune will not work correctly.
```
ssh root@pluto.local
```