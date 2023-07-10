/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package svrReader

import "q100transmitter/logger"

type (
	SvrConfig struct {
		IP_Address string
		IP_Port    int16
	}
	SvrData struct {
		Status string
	}
)

func Initialize(cfg SvrConfig, ch chan SvrData) {
	go readServer(cfg, ch)
}

func Stop() {
	logger.Warn.Printf("SvrClient will stop... - NOT IMPLELENTED")
	//
	logger.Info.Printf("SvrClient has stopped - NOT IMPLELENTED")
}
