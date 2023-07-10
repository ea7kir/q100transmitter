/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package encoderClient

import "q100transmitter/logger"

type (
	// API
	HeConfig struct {
		Audio_codec   string // "ACC"
		Audio_bitrate string // "64000"
		Video_codec   string // "H.265"
		Video_size    string // "1280x720"
		Video_bitrate string // "330"
		Url           string // "udp://192.168.3.10:8282"
		IP_Address    string // 192.168.3.1"
	}
)

// API
func Initialize(cfg *HeConfig) {
	setParams(cfg)
}

// API
//
// setarams is called from tuner. The function will write the params to a folder on the Pluto.
func SetParams(cfg *HeConfig) {
	setParams(cfg)
}

// see _MOTES/HVEC

func setParams(cfg *HeConfig) {
	logger.Info.Printf("writing params to the HEV-10 Encoder")
}
