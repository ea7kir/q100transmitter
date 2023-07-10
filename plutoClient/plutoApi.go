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

func Initialize(cfg PlConfig) {
	writePluto(&cfg)
}

// Called from tuner to copy the params into a folder in the Pluto.
func SetParams(cfg *PlConfig) {
	writePluto(cfg)
}
