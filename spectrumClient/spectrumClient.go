/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package spectrumClient

import (
	"q100transmitter/mylogger"

	"golang.org/x/net/websocket"
)

type (
	SpConfig struct {
		Url    string
		Origin string
	}
	SpData struct {
		Yp           []float32
		BeaconLevel  float32
		MarkerCentre float32
		MarkerWidth  float32
	}
)

var (
	Xp = make([]float32, numPoints) // x coordinates from 0.0 to 100.0
)

func Intitialize(cfg *SpConfig, ch chan SpData) {
	spChannel = ch
	Xp[0] = 0
	for i := 1; i < numPoints-1; i++ {
		Xp[i] = 100.0 * (float32(i) / float32(numPoints))
	}
	Xp[numPoints-1] = 100

	go readAndDecode(cfg, spChannel)
}

func Stop() {
	mylogger.Warn.Printf("Spectrum will stop... - NOT IMPLELENTED")
}

// Sets the spData Marker values
//
//	called from rxControl
func SetMarker(frequency, symbolRate string) {
	spData.MarkerCentre, spData.MarkerWidth = getMarkers(frequency, symbolRate)
	// spData.MarkerCentre = frequencyCentre[frequency]
	// spData.MarkerWidth = symbolRateWidth[symbolRate]
}

// END API *******************************************************

// room for 916 datapoints + start and end zero points to close the polygon
const numPoints = 918

var (
	spData = SpData{
		Yp:           make([]float32, numPoints),
		BeaconLevel:  0.5,
		MarkerCentre: 0.5,
		MarkerWidth:  0.5,
	}
	spChannel chan SpData
)

func readAndDecode(cfg *SpConfig, ch chan SpData) {
	ws, err := websocket.Dial(cfg.Url, "", cfg.Origin)
	if err != nil {
		mylogger.Fatal.Fatalf("Dial failed: %v", err)
	}
	defer ws.Close()

	var bytes = make([]byte, 2048) // larger than 1844
	var n int

	for {
		if n, err = ws.Read(bytes); err != nil {
			mylogger.Warn.Printf("Read failed: %v", err)
			continue
		}
		if n != 1844 {
			mylogger.Warn.Printf("reading : bytes != 1844\n")
			continue
		}

		// begin processing the bytes
		// count = 0
		for i := 0; i < 1836; {
			word := uint16(bytes[i]) + uint16(bytes[i+1])<<8
			// count++
			// mylogger.Info.Printf("count = %v\n", count)
			if word < 8192 {
				word = 8192
			}
			// spData.Yp[i/2] = float32(word-uint16(8192)) / float32(52000)
			spData.Yp[i/2] = float32(word-uint16(8192)) / float32(520) // normalize to 0 to 100
			// spData.Yp[i/2] = 50.0
			i += 2
		}
		// mylogger.Info.Printf("count = %v\n", count)
		spData.Yp[0] = 0
		spData.Yp[numPoints-1] = 0

		spData.BeaconLevel = 0
		for i := 32; i <= 133; i++ { // beacon center is 103
			spData.BeaconLevel += spData.Yp[i]
		}
		spData.BeaconLevel = spData.BeaconLevel / 103
		// mylogger.Info.Printf("beacon level %v : Yp[i] %v", spData.BeaconLevel, spData.Yp[103])

		ch <- spData
	}

}
