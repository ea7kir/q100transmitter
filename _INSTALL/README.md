# Installing Pi OS

NOTE: CURRENTLY REQUIRES PI OS BULLSEYE 64-BIT (FULL DESKTOP VERSION)

## Using Raspberry Pi Imager:

```
CHOOSE OS: Raspberry Pi OS -> Raspberry Pi OS (64-bit)

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

WARNING: the Pi will rebot during the install, so please allow it to complete

## Remote login from a Mac, PC or Linux host

```
ssh pi@q100transmitter.local
```

## More to follow goes here
