/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package plutoClient

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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

const (
	config_Url = "pluto.local" // or maybe "192.168.2.1"
)

type (
	PlConfig_t struct {
		Frequency       string // "2409.75"
		Mode            string // "DBS2"
		Constellation   string // "QPSK"
		SymbolRate      string // "333"
		Fec             string // "23"
		Gain            string // "-10"
		calibrationMode string // "nocalib"
		pcr_pts         string // "800"
		pat_period      string // "200"
		roll_off        string // "0.35"
		pilots          string // "off"
		frame           string // "LongFrame"
		h265box         string // "undefined"
		remux           string // "1"
		provider        string // "EA7KIR"
		service         string // "Michael"
	}
)

var (
	arg = PlConfig_t{}
)

func Start(provider, service string) {
	arg.provider = provider
	arg.service = service
	arg.calibrationMode = "nocalib"
	arg.pcr_pts = "800"
	arg.pat_period = "200"
	arg.roll_off = "0.35"
	arg.pilots = "off"
	arg.frame = "LongFrame"
	arg.h265box = "undefined"
	arg.remux = "1"
}

func Stop() {
	log.Printf("INFO plutoClient will stop... - NOT IMPLEMENTED")
	// now delete the local  settings file
	// ************* END dummy send using script

	// TODO: do it without using a script

	// can't do this until file is closed.
	// err = os.Remove(settingsFileName)
	// if err != nil {
	// 	log.Printf("WARN  Failed to delete settings.txt: %s", err)
	// }
}

// Called from tuner to copy the params into a folder in the Pluto.
func SetParams(cfg *PlConfig_t) {
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
		arg.provider,
		arg.Frequency,
		arg.Mode,
		arg.Constellation,
		arg.SymbolRate,
		arg.Fec,
		arg.pilots,
		arg.frame,
		arg.Gain,
		arg.roll_off,
		arg.pcr_pts,
		arg.pat_period,
		arg.h265box,
		arg.remux)

	// log.Printf("INFO 1: save to settings.txt to a local folder: \n%v\n", settings)

	const (
		cp2plutoScript   = "/home/pi/Q100/q100transmitter/etc/cp2pluto"
		settingsFileName = "/home/pi/Q100/settings.txt"
	)

	var (
		plutoDestination = "root@" + config_Url + ":/www/"
	)

	f, err := os.OpenFile(settingsFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("FATAL   %s", err)
	}
	defer f.Close()

	_, err = f.WriteString(settings)
	if err != nil {
		log.Fatalf("FATAL   %s", err)
	}
	// log.Printf("INFO Pluto settings saved to local file: %s", settingsFileName)

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
		log.Fatalf("FATAL   Failed to send %s to %s pluto: %s", settingsFileName, plutoDestination, err)
	}
}
