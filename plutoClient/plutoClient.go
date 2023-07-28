/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package plutoClient

import (
	"fmt"
	"os"
	"os/exec"
	"q100transmitter/logger"
	"strings"
)

/*
In the Pluto the file /www/settings.txt contains 14 feilds:

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
*/

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
		Url             string // "pluto.local" or "192.168.2.1",
	}
)

var (
	arg *PlConfig
)

func Initialize(cfg *PlConfig) {
	arg = cfg
	// // settings not used by the GUI
	// arg.Provider = cfg.Provider
	// arg.CalibrationMode = cfg.CalibrationMode
	// arg.Pcr_pts = cfg.Pcr_pts
	// arg.Pat_period = cfg.Pat_period
	// arg.Roll_off = cfg.Roll_off
	// arg.Pilots = cfg.Pilots
	// arg.Frame = cfg.Frame
	// arg.H265box = cfg.H265box
	// arg.Remux = cfg.Remux
}

// Called from tuner to copy the params into a folder in the Pluto.
func SetParams(cfg *PlConfig) {
	// overide settings provided by the GUI
	arg.Frequency = strings.Fields(cfg.Frequency)[0] // remove " / 27" etc
	arg.Mode = strings.Replace(cfg.Mode, "-", "", 1) // remove "-""
	arg.Constellation = cfg.Constellation
	arg.SymbolRate = cfg.SymbolRate
	arg.Fec = strings.Replace(cfg.Fec, "/", "", 1) // remove "/""
	arg.Gain = cfg.Gain
	writePluto()
}

func writePluto() {
	settings := fmt.Sprintf("callsign %v\nfreq %v\nmode %v\nmod %v\nsr %v\nfec %v\npilots %v\nframe %v\npower %v\nrolloff %v\npcrpts %v\npatperiod %v\nh265box %v\nremux %v\n\n",
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

	// logger.Info.Printf("1: save to settings.txt to a local folder: \n%v\n", settings)

	const cp2plutoScript = "/home/pi/Q100/q100transmitter/_scripts/cp2pluto"
	const settingsFileName = "/home/pi/Q100/settings.txt"
	const plutoDestination = "root@pluto.local:/www/"

	f, err := os.OpenFile(settingsFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Fatal.Fatalf("%s", err)
	}
	defer f.Close()

	_, err = f.WriteString(settings)
	if err != nil {
		logger.Fatal.Fatalf("%s", err)
	}
	// logger.Info.Printf("Pluto settings saved to local file: %s", settingsFileName)

	// Sending to Pluto on the smd line
	// /usr/bin/sshpass -panalog /usr/bin/scp /home/pi/settings.txt root@pluto.local:/www/ > /dev/null 2>&1

	// ************* BEGIN dummy send using script

	// args :=
	cmd := exec.Command(cp2plutoScript, // TODO: args here are currently ignored by the script
		// "/usr/bin/sshpass",
		// "-panalog",
		// "/usr/bin/scp",
		settingsFileName,
		plutoDestination,
	)
	_, err = cmd.Output()
	if err != nil {
		logger.Fatal.Fatalf("Failed to send %s to %s pluto: %s", settingsFileName, plutoDestination, err)
		return //false
	}
	// now delete the local  settings file
	// ************* END dummy send using script

	// TODO: do it without using a script

	// can't do this until file is closed.
	// err = os.Remove(settingsFileName)
	// if err != nil {
	// 	logger.Warn.Printf("Failed to delete settings.txt: %s", err)
	// }
}

/*
THE PYTHON WAY

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

Now we can ssh and scp like this...

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
