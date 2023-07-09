/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package spReader

import "q100transmitter/logger"

type (
	SpConfig struct {
		Url string
	}
	SpData struct {
		Yp                        []float32
		BeaconLevel               float32
		MarkerCentre, MarkerWidth float32
	}
)

var (
	cfg *SpConfig
	Xp  = make([]float32, numPoints) // x coordinates from 0.0 to 100.0
)

func Intitialize(spc SpConfig, ch chan SpData) {
	cfg = &spc
	// url = spc.Url
	spChannel = ch
	Xp[0] = 0
	for i := 1; i < numPoints-1; i++ {
		Xp[i] = 100.0 * (float32(i) / float32(numPoints))
	}
	Xp[numPoints-1] = 100

	go readAndDecode(cfg.Url, spChannel)
}

// func Start() {
// 	// logger.Info.Printf("spectrum.readAndDecode will start...")
// 	logger.Info.Printf("Spectrum will start...")
// 	go readAndDecode(cfg.Url, spChannel)
// 	logger.Info.Printf("Spectrum has started")
// }

func Stop() {
	logger.Info.Printf("Spectrum will stop... - DOES NOTHING")
	//
	logger.Info.Printf("Spectrum has stopped")
}

func SetMarker(frequency, symbolRate string) {
	spData.MarkerCentre, spData.MarkerWidth = getMarkers(frequency, symbolRate)
}
