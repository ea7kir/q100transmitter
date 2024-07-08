package spClient

/*****************************************************************
 * SPECTRUM MARKERS FOR TRANSMITING
 *****************************************************************/

var (
	// TODO: calculatee a mathematical values
	const_frequencyCentre = map[string]float32{
		// "2402.00 / 00": 103, // Beacon not required for Tx
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
	centre := const_frequencyCentre[frequency] / 9.18 // NOTE: 9.18 is a temporary kludge
	width := const_symbolRateWidth[symbolRate]
	return centre, width
}

// Sets the spData Marker values
func SetMarker(frequency string, symbolRate string) {
	spData.MarkerCentre, spData.MarkerWidth = getMarkers(frequency, symbolRate)
}
