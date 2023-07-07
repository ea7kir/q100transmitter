/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package hev10

import "q100transmitter/logger"

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
	//hecfg = cfg
}

func Config(cfg HeConfig) {
	logger.Info.Printf("will configure HEV-10...")
	config(cfg)
	logger.Info.Printf("has configured HEV-10")
}

func UnConfig() {
	logger.Info.Printf("will unconfigure HEV-10...")

	logger.Info.Printf("has unconfigured HEV-10")
}
