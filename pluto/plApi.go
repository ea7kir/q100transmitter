/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package pluto

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
	}
)

func Intitialize(args PlConfig) {
	// plc = cfg
}