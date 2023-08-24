/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package spectrumClient

/*****************************************************************
* SPECTRUM & CALIBRARTION MARKERS
*****************************************************************/

var (
	// TODO: calculatee a mathematical values
	const_frequencyCentre = map[string]float32{
		"2403.25 / 01": 230, // / divisor,
		"2403.50 / 02": 256, // / divisor,
		"2403.75 / 03": 281, // / divisor,
		"2404.00 / 04": 307, // / divisor,
		"2404.25 / 05": 332, // / divisor,
		"2404.50 / 06": 358, // / divisor,
		"2404.75 / 07": 383, // / divisor,
		"2405.00 / 08": 409, // / divisor,
		"2405.25 / 09": 434, // / divisor,
		"2405.50 / 10": 460, // / divisor,
		"2405.75 / 11": 485, // / divisor,
		"2406.00 / 12": 511, // / divisor,
		"2406.25 / 13": 536, // / divisor,
		"2406.50 / 14": 562, // / divisor,
		"2406.75 / 15": 588, // / divisor,
		"2407.00 / 16": 613, // / divisor,
		"2407.25 / 17": 639, // / divisor,
		"2407.50 / 18": 664, // / divisor,
		"2407.75 / 19": 690, // / divisor,
		"2408.00 / 20": 715, // / divisor,
		"2408.25 / 21": 741, // / divisor,
		"2408.50 / 22": 767, // / divisor,
		"2408.75 / 23": 792, // / divisor,
		"2409.00 / 24": 818, // / divisor,
		"2409.25 / 25": 843, // / divisor,
		"2409.50 / 26": 869, // / divisor,
		"2409.75 / 27": 894, // / divisor,
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
