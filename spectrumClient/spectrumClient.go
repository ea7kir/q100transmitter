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
/* commented out this doesn't work here
var (
	// TODO: calculatee a mathematical values
	const_frequencyCentre = map[string]float32{
		"10491.50 / 00": 103,
		"10492.75 / 01": 230,
		"10493.00 / 02": 256,
		"10493.25 / 03": 281,
		"10493.50 / 04": 307,
		"10493.75 / 05": 332,
		"10494.00 / 06": 358,
		"10494.25 / 07": 383,
		"10494.50 / 08": 409,
		"10494.75 / 09": 434,
		"10495.00 / 10": 460,
		"10495.25 / 11": 485,
		"10495.50 / 12": 511,
		"10495.75 / 13": 536,
		"10496.00 / 14": 562,
		"10496.25 / 15": 588,
		"10496.50 / 16": 613,
		"10496.75 / 17": 639,
		"10497.00 / 18": 664,
		"10497.25 / 19": 690,
		"10497.50 / 20": 715,
		"10497.75 / 21": 741,
		"10490.00 / 22": 767,
		"10498.25 / 23": 792,
		"10498.50 / 24": 818,
		"10498.75 / 25": 843,
		"10499.00 / 26": 869,
		"10499.25 / 27": 894,
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
func getMarkers(frequency, sysmbolRate string) (float32, float32) {
	centre := const_frequencyCentre[frequency] / 9.18 // 9.18 is a temporary kludge
	width := const_symbolRateWidth[sysmbolRate]
	return centre, width
}

commmented out */

// TODO: implement CalibratetionPoints()
/*
 func CalibratetionPoints() {
	 var yp [918]float32

	 for _, v := range CalibrationMarkerWidth {
		 // yp[v] = 100
		 mylogger.Info.Printf("CalibratetionPoints %v", v)
	 }

	 for i, v := range yp {
		 mylogger.Info.Printf("CalibratetionPoints %v  %v", i, v)
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
