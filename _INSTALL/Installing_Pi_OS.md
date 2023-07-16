# Installing Pi OS

NOTE: CURRENTLY REQUIRES PIOS BULLSEYE 64-BIT LITE

## Using Raspberry Pi Imager:

```
CHOOSE OS:	Raspberry Pi OS (64-bit)

CONFIGURE:
	Set hostname:			txtouch
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

Insert the card and login. Wait for the software update icon to appear, and proceed to update and reboot.

## Login

Clone the repro.

```
cd
git clone https://github.com/ea7kir/TxTouch.git
```

Change permissions for all the install scrpits.

```
chmod +x TxTouch/_Resources/install_*
```

Exeute each script in numerical order.


