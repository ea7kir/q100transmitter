/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package paClient

import (
	"fmt"
	"q100transmitter/logger"
	"time"
)

type (
	// API
	SvrConfig struct {
		IP_Address string
		IP_Port    int16
	}
	// API
	SvrData struct {
		Status string
	}
)

// API
func Initialize(cfg *SvrConfig, ch chan SvrData) {
	go readServer(cfg, ch)
}

// API
func Stop() {
	logger.Warn.Printf("SvrClient will stop... - NOT IMPLELENTED")
	//
	logger.Info.Printf("SvrClient has stopped - NOT IMPLELENTED")
}

func readServer(cfg *SvrConfig, ch chan SvrData) {

	sd := SvrData{}
	count := 0

	for {
		time.Sleep(time.Second)
		count++
		str := fmt.Sprintf("SUCCESS # %v", count)
		sd.Status = str
		ch <- sd
	}

}
