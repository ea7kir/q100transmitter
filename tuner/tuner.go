/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package tuner

var (
	const_BAND_LIST = []string{
		"Wide",
		"Narrow",
		"V.Narrow",
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
	const_WIDE_FREQUENCY_LIST = []string{
		"2403.75 / 03",
		"2405.25 / 09",
		"2406.75 / 15",
	}
	const_NARROW_FREQUENCY_LIST = []string{
		"2403.25 / 01",
		"2403.75 / 03",
		"2404.25 / 05",
		"2404.75 / 07",
		"2405.25 / 09",
		"2405.75 / 11",
		"2406.25 / 13",
		"2406.75 / 15",
		"2407.25 / 17",
		"2407.75 / 19",
		"2408.25 / 21",
		"2408.75 / 23",
		"2409.25 / 25",
		"2409.75 / 27", // _f_index 13
	}
	const_VERY_NARROW_FREQUENCY_LIST = []string{
		"2403.25 / 01",
		"2403.50 / 02",
		"2403.75 / 03",
		"2404.00 / 04",
		"2404.25 / 05",
		"2404.50 / 06",
		"2404.75 / 07",
		"2405.00 / 08",
		"2405.25 / 09",
		"2405.50 / 10",
		"2405.75 / 11",
		"2406.00 / 12",
		"2406.25 / 13",
		"2406.50 / 14",
		"2406.75 / 15",
		"2407.00 / 16",
		"2407.25 / 17",
		"2407.50 / 18",
		"2407.75 / 19",
		"2408.00 / 20",
		"2408.25 / 21",
		"2408.50 / 22",
		"2408.75 / 23",
		"2409.00 / 24",
		"2409.25 / 25",
		"2409.50 / 26",
		"2409.75 / 27",
	}
	const_WIDE_MODE_LIST = []string{
		"DVB-S", "DVB-S2",
	}
	const_NARROW_MODE_LIST = []string{
		"DVB-S", "DVB-S2",
	}
	const_VERY_NARROW_MODE_LIST = []string{
		"DVB-S", "DVB-S2",
	}
	const_WIDE_CODECS_LIST = []string{
		"H264 ACC", "H265 ACC",
	}
	const_NARROW_CODECS_LIST = []string{
		"H264 ACC", "H265 ACC",
	}
	const_VERY_NARROW_CODECS_LIST = []string{
		"H264 ACC", "H265 ACC",
	}
	const_WIDE_CONSTELLATION_LIST = []string{
		"QPSK", "8PSK", "16PSK", "32PSK",
	}
	const_NARROW_CONSTELLATION_LIST = []string{
		"QPSK", "8PSK", "16PSK", "32PSK",
	}
	const_VERY_NARROW_CONSTELLATION_LIST = []string{
		"QPSK", "8PSK", "16PSK", "32PSK",
	}
	const_WIDE_FEC_LIST = []string{
		"1/2", "2/3", "3/4", "4/5", "5/6", "6/7", "7/8", "8/9",
	}
	const_NARROW_FEC_LIST = []string{
		"1/2", "2/3", "3/4", "4/5", "5/6", "6/7", "7/8", "8/9",
	}
	const_VERY_NARROW_FEC_LIST = []string{
		"1/2", "2/3", "3/4", "4/5", "5/6", "6/7", "7/8", "8/9",
	}
	const_WIDE_VIDEO_BITRATE_LIST = []string{
		"290", "300", "310", "330", "340", "350", "360",
	}
	const_NARROW_VIDEO_BITRATE_LIST = []string{
		"290", "300", "310", "330", "340", "350", "360",
	}
	const_VERY_NARROW_VIDEO_BITRATE_LIST = []string{
		"290", "300", "310", "330", "340", "350", "360",
	}
	const_WIDE_AUDIO_BITRATE_LIST = []string{
		"32000", "64000",
	}
	const_NARROW_AUDIO_BITRATE_LIST = []string{
		"32000", "64000",
	}
	const_VERY_NARROW_AUDIO_BITRATE_LIST = []string{
		"32000", "64000",
	}
	const_WIDE_SPARE1_LIST = []string{
		"sp1-a", "sp1-b",
	}
	const_NARROW_SPARE1_LIST = []string{
		"sp1-a", "sp1-b",
	}
	const_VERY_NARROW_SPARE1_LIST = []string{
		"sp1-a", "sp1-b",
	}
	const_WIDE_SPARE2_LIST = []string{
		"sp2-a", "sp2-b",
	}
	const_NARROW_SPARE2_LIST = []string{
		"sp2-a", "sp2-b",
	}
	const_VERY_NARROW_SPARE2_LIST = []string{
		"sp2-a", "sp2-b",
	}
	const_WIDE_GAIN_LIST = []string{
		"-20", "-19", "-18", "-17", "-16", "-15", "-14", "-13", "-12", "-11", "-10", // "-9","-8","-7","-6","-5","-4","-3","-2","-1","0",
	}
	const_NARROW_GAIN_LIST = []string{
		"-20", "-19", "-18", "-17", "-16", "-15", "-14", "-13", "-12", "-11", "-10", // "-9","-8","-7","-6","-5","-4","-3","-2","-1","0",
	}
	const_VERY_NARROW_GAIN_LIST = []string{
		"-20", "-19", "-18", "-17", "-16", "-15", "-14", "-13", "-12", "-11", "-10", //"-9","-8","-7","-6","-5","-4","-3","-2","-1","0",
	}

	wideSymbolRate          Selector
	narrowSymbolRate        Selector
	veryNarrowSymbolRate    Selector
	wideFrequency           Selector
	narrowFrequency         Selector
	veryNarrowFrequency     Selector
	wideMode                Selector
	narrowMode              Selector
	veryNarrowMode          Selector
	wideCodecs              Selector
	narrowCodecs            Selector
	veryNarrowCodecs        Selector
	wideConstellation       Selector
	narrowConstellation     Selector
	veryNarrowConstellation Selector
	wideFec                 Selector
	narrowFec               Selector
	veryNarrowFec           Selector
	wideVideoBitRate        Selector
	narrowVideoBitRate      Selector
	veryNarrowVideoBitRate  Selector
	wideAudioBitRate        Selector
	narrowAudioBitRate      Selector
	veryNarrowAudioBitRate  Selector
	wideSpare1              Selector
	narrowSpare1            Selector
	veryNarrowSpare1        Selector
	wideSpare2              Selector
	narrowSpare2            Selector
	veryNarrowSpare2        Selector
	wideGain                Selector
	narrowGain              Selector
	veryNarrowGain          Selector
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

func newSelector(values []string, with string) Selector {
	index := indexInList(values, with)
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

	case const_BAND_LIST[0]: // wide
		SymbolRate = wideSymbolRate
		Frequency = wideFrequency
		Mode = wideMode
		Codecs = wideCodecs
		Constellation = wideConstellation
		Fec = wideFec
		VideoBitRate = wideVideoBitRate
		AudioBitRate = wideAudioBitRate
		Spare1 = wideSpare1
		Spare2 = wideSpare2
		Gain = wideGain
	case const_BAND_LIST[1]: // narrow
		SymbolRate = narrowSymbolRate
		Frequency = narrowFrequency
		Mode = narrowMode
		Codecs = narrowCodecs
		Constellation = narrowConstellation
		Fec = narrowFec
		VideoBitRate = narrowVideoBitRate
		AudioBitRate = narrowAudioBitRate
		Spare1 = narrowSpare1
		Spare2 = narrowSpare2
		Gain = narrowGain
	case const_BAND_LIST[2]: // very narrow
		SymbolRate = veryNarrowSymbolRate
		Frequency = veryNarrowFrequency
		Mode = veryNarrowMode
		Codecs = veryNarrowCodecs
		Constellation = veryNarrowConstellation
		Fec = veryNarrowFec
		VideoBitRate = veryNarrowVideoBitRate
		AudioBitRate = veryNarrowAudioBitRate
		Spare1 = veryNarrowSpare1
		Spare2 = veryNarrowSpare2
		Gain = veryNarrowGain
	}
}