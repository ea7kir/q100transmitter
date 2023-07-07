/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package hev10

import "q100transmitter/logger"

func Config() {
	logger.Info.Printf("will configure HEV-10...")

	logger.Info.Printf("has configured HEV-10")
}

func UnConfig() {
	logger.Info.Printf("will unconfigure HEV-10...")

	logger.Info.Printf("has unconfigured HEV-10")
}
