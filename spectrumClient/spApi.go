/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package spectrumClient

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
	Xp = make([]float32, numPoints) // x coordinates from 0.0 to 100.0
)

func Intitialize(cfg SpConfig, ch chan SpData) {
	spChannel = ch
	Xp[0] = 0
	for i := 1; i < numPoints-1; i++ {
		Xp[i] = 100.0 * (float32(i) / float32(numPoints))
	}
	Xp[numPoints-1] = 100

	go readAndDecode(cfg.Url, spChannel)
}

func Stop() {
	logger.Warn.Printf("Spectrum will stop... - NOT IMPLELENTED")
	//
	logger.Info.Printf("Spectrum has stopped - NOT IMPLELENTED")
}

// Sets the spData Marker values
func SetMarker(frequency, symbolRate string) {
	spData.MarkerCentre, spData.MarkerWidth = getMarkers(frequency, symbolRate)
}
