## TODO:

- add wiring schematic photos and better info to the doc folder
- revise what data to monitor - eg nulls
- revise what parameters to use
- improve marker widths
- more to do in spClient to deal with doubling / fast changes /etc

# Find ways to make install easier

- Appearance Settings / Desktop:
    - Disable Documents, Wastebasket and External Disks
- TurnOff Bluetooth

## Auto Start
- Currently using systemctl and NOT wayfire.ini for run at boot
    - because ~/wayland.ini ```[autostart]``` isn't behaving

## Maybe one day

- move from Bookworm Desktop to Bookworm Lite or FreeBSD

## Possible ways to disbale bt and wifi

- /boot/firmware/config.txt
	- dtoverlay=disable-wifi
	- dtoverlay=disable-bt
- or maybe
	- sudo systemctl disable bluetooth.service
	- sudo systemctl stop bluetooth.service
- or maybe
	- sudo rfkill block wifi
	- sudo rfkill block bluetooth
- or maybe
	- sudo systemctl disable wpa_supplicant
	- sudo systemctl disable bluetooth
	- sudo systemctl disable hciuart
