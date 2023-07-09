/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package hev10

type (
	HeConfig struct {
		Audio_codec   string // "ACC"
		Audio_bitrate string // "64000"
		Video_codec   string // "H.265"
		Video_size    string // "1280x720"
		Video_bitrate string // "330"
		Url           string // "udp://192.168.3.10:8282"
	}
)

func Initialize(cfg HeConfig) {
	setParams(&cfg)
}

// setarams is called from tuner. The function will write the params to a folder on the Pluto.
func SetParams(cfg *HeConfig) {
	setParams(cfg)
}
