/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package plutoClient

import (
	"fmt"
	"q100transmitter/logger"
)

/*
import subprocess
#import sys

def configure_pluto():
    pass

def shutdown_pluto():
    pass

"""
In the Pluto the file /www/settings.txt contains:

callsign EA7KIR
freq 2409.75
mode DVBS2
mod QPSK
sr 333
fec 34
pilots Off
frame LongFrame
power -2
rolloff 0.25
pcrpts 800
patperiod 200
h265box undefined
remux 1

So, we will create a local file and copy it to the pluto.
It will be neccessary to avoid login and passwprd, so install the following...

See: https://wiki.analog.com/university/tools/pluto/drivers/linux
sudo apt install libiio-utils
iio_info -n 192.168.2.1 | grep device
io_readdev -n 192.168.2.1 -s 64 cf-ad9361-lpc | hexdump -x

AND FOR PASSWORDS:

wget https://raw.githubusercontent.com/analogdevicesinc/plutosdr_scripts/master/ssh_config -O ~/.ssh/config
sudo apt install sshpass

Now I can ssh and scp like this...

sshpass -panalog ssh root@pluto.local
sshpass -panalog ssh root@192.168.2.1 # not working
sshpass -panalog scp /home/pi/settings.txt root@pluto.local:/www/
sshpass -panalog scp /home/pi/settings.txt root@192.168.2.1:/www/  # not working
"""

def setup_pluto(args):
    # stop_pluto()
    settings = 'callsign {}\nfreq {}\nmode {}\nmod {}\nsr {}\nfec {}\npilots {}\nframe {}\npower {}\nrolloff {}\npcrpts {}\npatperiod {}\nh265box {}\nremux {}\n\n'.format(
        args.provider,
        args.frequency,
        args.mode,
        args.constellation,
        args.symbol_rate,
        args.fec,
        args.pilots,
        args.frame,
        args.gain,
        args.roll_off,
        args.pcr_pts,
        args.pat_period,
        args.h265box,
        args.remux)
    f = open("/home/pi/settings.txt", "w")
    f.write(settings)
    f.close()

    #args = ['/usr/bin/sshpass', '-panalog', '/usr/bin/scp', '/home/pi/settings.txt', 'root@pluto.local:/www/']
    #result = subprocess.run(args)

    cmd_str = '/usr/bin/sshpass -panalog /usr/bin/scp /home/pi/settings.txt root@pluto.local:/www/ > /dev/null 2>&1'
    result = subprocess.run(cmd_str, shell=True)

    if result.returncode != 0:
        print('ERROR updating pluto settings.txt', flush=True)
    #else:
    #    print('pluto configured ok', flush=True)
    # start_pluto()


*/

/*
	plConfig = plutoClient.PlConfig{
		Frequency:        "2409.75",
		Mode:             "DBS2",
		Constellation:    "QPSK",
		Symbol_rate:      "333",
		Fec:              "23",
		Gain:             "-10",
		Calibration_mode: "nocalib",   // NOTE: not implemented
		Pcr_pts:          "800",       // NOTE: not implemented
		Pat_period:       "200",       // NOTE: not implemented
		Roll_off:         "0.35",      // NOTE: not implemented
		Pilots:           "off",       // NOTE: not implemented
		Frame:            "LongFrame", // NOTE: not implemented
		H265box:          "undefined", // NOTE: not implemented
		Remux:            "1",         // NOTE: not implemented
		Provider:         "EA7KIR",
		Service:          "Michael", // NOTE: not implemented
		IP_Address:       "192.168.2.1",
	}
*/
func writePluto(cfg *PlConfig) {
	str := fmt.Sprintf("callsign %v\nfreq %v\nmode %v\nmod %v\nsr %v\nfec %v\npilots %v\nframe %v\npower %v\nrolloff %v\npcrpts %v\npatperiod %v\nh265box %v\nremux %v\n\n",
		cfg.Provider,
		cfg.Frequency,
		cfg.Mode,
		cfg.Constellation,
		cfg.Symbol_rate,
		cfg.Fec,
		cfg.Pilots,
		cfg.Frame,
		cfg.Gain,
		cfg.Roll_off,
		cfg.Pcr_pts,
		cfg.Pat_period,
		cfg.H265box,
		cfg.Remux)

	logger.Info.Printf("writing params to a folder on the Pluto: \n%v\n", str)
}
