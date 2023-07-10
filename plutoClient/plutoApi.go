/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package plutoClient

type (
	PlConfig struct {
		Frequency        string // "2409.75"
		Mode             string // "DBS2"
		Constellation    string // "QPSK"
		Symbol_rate      string // "333"
		Fec              string // "23"
		Gain             string // "-10"
		Calibration_mode string // "nocalib"
		Pcr_pts          string // "800"
		Pat_period       string // "200"
		Roll_off         string // "0.35"
		Pilots           string // "off"
		Frame            string // "LongFrame"
		H265box          string // "undefined"
		Remux            string // "1"
		Provider         string // "EA7KIR"
		Service          string // "Michael"
		IP_Address       string // "192.168.2.1",
	}
)

func Intitialize(cfg PlConfig) {
	writePluto(&cfg)
}

// Called from tuner to copy the params into a folder in the Pluto.
func SetParams(cfg *PlConfig) {
	writePluto(cfg)
}
