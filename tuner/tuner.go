/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package tuner

var (
	const_BAND_LIST = []string{
		"Beacon",
		"Wide",
		"Narrow",
		"V.Narrow",
	}
	const_BEACON_FREQUENCY_LIST = []string{
		"10491.50 / 00",
	}
	const_WIDE_FREQUENCY_LIST = []string{
		"10493.25 / 03",
		"10494.75 / 09",
		"10496.25 / 15",
	}
	const_NARROW_FREQUENCY_LIST = []string{
		"10492.75 / 01",
		"10493.25 / 03",
		"10493.75 / 05",
		"10494.25 / 07",
		"10494.75 / 09",
		"10495.25 / 11",
		"10495.75 / 13",
		"10496.25 / 15",
		"10496.75 / 17",
		"10497.25 / 19",
		"10497.75 / 21",
		"10498.25 / 23",
		"10498.75 / 25",
		"10499.25 / 27", // index 13
	}
	const_VERY_NARROW_FREQUENCY_LIST = []string{
		"10492.75 / 01",
		"10493.00 / 02",
		"10493.25 / 03",
		"10493.50 / 04",
		"10493.75 / 05",
		"10494.00 / 06",
		"10494.25 / 07",
		"10494.50 / 08",
		"10494.75 / 09",
		"10495.00 / 10",
		"10495.25 / 11",
		"10495.50 / 12",
		"10495.75 / 13",
		"10496.00 / 14", // index 13
		"10496.25 / 15",
		"10496.50 / 16",
		"10496.75 / 17",
		"10497.00 / 18",
		"10497.25 / 19",
		"10497.50 / 20",
		"10497.75 / 21",
		"10498.00 / 22",
		"10498.25 / 23",
		"10498.50 / 24",
		"10498.75 / 25",
		"10499.00 / 26",
		"10499.25 / 27",
	}
	const_BEACON_SYMBOLRATE_LIST = []string{
		"1500",
	}
	const_WIDE_SYMBOLRATE_LIST = []string{
		"1000",
		"1500",
		"2000",
	}
	const_NARROW_SYMBOLRATE_LIST = []string{
		"250",
		"333",
		"500",
	}
	const_VERY_NARROW_SYMBOLRATE_LIST = []string{
		"33",
		"66",
		"125",
	}

	beaconSymbolRate     Selector
	beaconFrequency      Selector
	wideSymbolRate       Selector
	wideFrequency        Selector
	narrowSymbolRate     Selector
	narrowFrequency      Selector
	veryNarrowSymbolRate Selector
	veryNarrowFrequency  Selector
)

// TODO: add error chatching
func indexInList(list []string, with string) int {
	for i := range list {
		if list[i] == with {
			return i
		}
	}
	return 0
}

func newSelector(values []string, index int) Selector {
	st := Selector{
		currIndex: index,
		lastIndex: len(values) - 1,
		list:      values,
		Value:     values[index],
	}
	return st
}

func switchBand() { // TODO: should switch back to previosly use settings
	switch Band.Value {
	case const_BAND_LIST[0]: // beacon
		SymbolRate = beaconSymbolRate
		Frequency = beaconFrequency
	case const_BAND_LIST[1]: // wide
		SymbolRate = wideSymbolRate
		Frequency = wideFrequency
	case const_BAND_LIST[2]: // narrow
		SymbolRate = narrowSymbolRate
		Frequency = narrowFrequency
	case const_BAND_LIST[3]: // very narrow
		SymbolRate = veryNarrowSymbolRate
		Frequency = veryNarrowFrequency
	}
}
