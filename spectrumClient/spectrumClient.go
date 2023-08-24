/*
 *  Q-100 Receiver & Transmitter
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
//	called from rxControl or txControl
func SetMarker(frequency string, symbolRate string) {
	spData.MarkerCentre, spData.MarkerWidth = getMarkers(frequency, symbolRate)
	// spData.MarkerCentre = const_frequencyCentre[frequency]
	// spData.MarkerWidth = const_symbolRateWidth[symbolRate]
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

/*****************************************************************
* SPECTRUM & CALIBRARTION MARKERS
*****************************************************************/

var (
	// TODO: calculatee a mathematical values
	const_frequencyCentre = map[string]float32{
		"2403.25 / 01": 230,
		"2403.50 / 02": 256,
		"2403.75 / 03": 281,
		"2404.00 / 04": 307,
		"2404.25 / 05": 332,
		"2404.50 / 06": 358,
		"2404.75 / 07": 383,
		"2405.00 / 08": 409,
		"2405.25 / 09": 434,
		"2405.50 / 10": 460,
		"2405.75 / 11": 485,
		"2406.00 / 12": 511,
		"2406.25 / 13": 536,
		"2406.50 / 14": 562,
		"2406.75 / 15": 588,
		"2407.00 / 16": 613,
		"2407.25 / 17": 639,
		"2407.50 / 18": 664,
		"2407.75 / 19": 690,
		"2408.00 / 20": 715,
		"2408.25 / 21": 741,
		"2408.50 / 22": 767,
		"2408.75 / 23": 792,
		"2409.00 / 24": 818,
		"2409.25 / 25": 843,
		"2409.50 / 26": 869,
		"2409.75 / 27": 894,
	}

	// TODO: calculatee a mathematical values
	const_symbolRateWidth = map[string]float32{
		"2000": 20,
		"1500": 15,
		"1000": 10,
		"500":  8,
		"333":  5,
		"250":  4,
		"125":  3,
		"66":   2,
		"33":   1.5,
	}
)

// Returns frequency and bandWidth Markers as float32
func getMarkers(frequency, symbolRate string) (float32, float32) {
	centre := const_frequencyCentre[frequency] / 9.18 // 9.18 is a temporary kludge
	width := const_symbolRateWidth[symbolRate]
	return centre, width
}

// TODO: implement CalibratetionPoints()
/*
func CalibratetionPoints() {
	var yp [918]float32

	for _, v := range CalibrationMarkerWidth {
		// yp[v] = 100
		logger.Info.Printf("CalibratetionPoints %v", v)
	}

	for i, v := range yp {
		logger.Info.Printf("CalibratetionPoints %v  %v", i, v)
	}

}
*/

/*
func readCalibrationData(ch chan SpData) {
	mylogger.Info.Printf("Spectrun calibration running...")
	for {
		spData.Yp[0] = 0
		for i := 1; i < numPoints-2; i++ {
			spData.Yp[i] = rand.Float32() * 50.0
		}
		spData.Yp[numPoints-1] = 0
		spData.BeaconLevel = rand.Float32() * 100
		ch <- spData
		time.Sleep(3 * time.Millisecond)
	}
}
*/
