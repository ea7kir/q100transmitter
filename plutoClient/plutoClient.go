/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package plutoClient

import (
	"fmt"
	"q100transmitter/logger"
)

// API
type (
	PlConfig struct {
		Frequency       string // "2409.75"
		Mode            string // "DBS2"
		Constellation   string // "QPSK"
		SymbolRate      string // "333"
		Fec             string // "23"
		Gain            string // "-10"
		CalibrationMode string // "nocalib"
		Pcr_pts         string // "800"
		Pat_period      string // "200"
		Roll_off        string // "0.35"
		Pilots          string // "off"
		Frame           string // "LongFrame"
		H265box         string // "undefined"
		Remux           string // "1"
		Provider        string // "EA7KIR"
		Service         string // "Michael"
		IP_Address      string // "192.168.2.1",

		// callsign  string // EA7KIR
		// freq      string // 2409.75
		// mode      string // DVBS2
		// mod       string // QPSK
		// sr        string // 333
		// fec       string // 34
		// pilots    string // Off
		// frame     string // LongFrame
		// power     string // -2
		// rolloff   string // 0.25
		// pcrpts    string // 800
		// patperiod string // 200
		// h265box   string // undefined
		// remux     string // 1

	}
)

var (
	arg = PlConfig{}
)

func Initialize(cfg *PlConfig) {
	arg.CalibrationMode = cfg.CalibrationMode // NOTE: not implemented
	arg.Pcr_pts = cfg.Pcr_pts                 // NOTE: not implemented
	arg.Pat_period = cfg.Pat_period           // NOTE: not implemented
	arg.Roll_off = cfg.Roll_off               // NOTE: not implemented
	arg.Pilots = cfg.Pilots                   // NOTE: not implemented
	arg.Frame = cfg.Frame                     // NOTE: not implemented
	arg.H265box = cfg.H265box                 // NOTE: not implemented
	arg.Remux = cfg.Remux                     // NOTE: not implementearg
}

// Called from tuner to copy the params into a folder in the Pluto.
func SetParams(cfg *PlConfig) {

	arg.Provider = cfg.Provider
	arg.Frequency = cfg.Frequency
	arg.Mode = cfg.Mode
	arg.Constellation = cfg.Constellation
	arg.SymbolRate = cfg.SymbolRate
	arg.Fec = cfg.Fec
	arg.Gain = cfg.Gain
	writePluto()
}

func writePluto() {
	str := fmt.Sprintf("callsign %v\nfreq %v\nmode %v\nmod %v\nsr %v\nfec %v\npilots %v\nframe %v\npower %v\nrolloff %v\npcrpts %v\npatperiod %v\nh265box %v\nremux %v\n\n",
		arg.Provider,
		arg.Frequency,
		arg.Mode,
		arg.Constellation,
		arg.SymbolRate,
		arg.Fec,
		arg.Pilots,
		arg.Frame,
		arg.Gain,
		arg.Roll_off,
		arg.Pcr_pts,
		arg.Pat_period,
		arg.H265box,
		arg.Remux)

	logger.Info.Printf("writing params to a folder on the Pluto: \n%v\n", str)
}

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
