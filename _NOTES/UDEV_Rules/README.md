# UDEV Rules

From:

`https://wiki.analog.com/university/tools/pluto/drivers/linux`

In order to access some USB functions without root privileges, it's recommended to install the PlutoSDR udev rules.

Download rules from::

`https://raw.githubusercontent.com/analogdevicesinc/plutosdr-fw/master/scripts/53-adi-plutosdr-usb.rules`

And copy: 

`sudo cp 53-adi-plutosdr-usb.rules /etc/udev/rules.d/`

Reload rules:

`sudo udevadm control --reload-rules`



