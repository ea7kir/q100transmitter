/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package txControl

import (
	"context"
	"log"
	"q100transmitter/encoderClient"
	"q100transmitter/plutoClient"
	"q100transmitter/pttSwitch"
)

const (
	config_Band                    = "Narrow"
	config_WideSymbolrate          = "1000"
	config_NarrowSymbolrate        = "333"
	config_VeryNarrowSymbolRate    = "125"
	config_WideFrequency           = "2405.25 / 09"
	config_NarrowFrequency         = "2409.75 / 27"
	config_VeryNarrowFrequency     = "2406.50 / 14"
	config_WideMode                = "DVB-S2"
	config_NarrowMode              = "DVB-S2"
	config_VeryNarrowMode          = "DVB-S2"
	config_WideCodecs              = "H265 ACC" // H.264 ACC | H.264 G711u | H.265 ACC | H.265 G711u
	config_NarrowCdecs             = "H265 ACC"
	config_VeryNarrowCodecs        = "H265 ACC"
	config_WideConstellation       = "QPSK"
	config_NarrowConstellation     = "QPSK"
	config_VeryNarrorConstellation = "QPSK"
	config_WideFec                 = "3/4"
	config_NarrowFec               = "3/4"
	config_VeryNarrowFec           = "3/4"
	config_WideVideoBitRate        = "1100" // 32...16384
	config_NarrowVideoBitRate      = "380"
	config_VeryNarrowVideoBitRate  = "250"
	config_WideAudioBitRate        = "48000" // 32000 | 48000 | 64000
	config_NarrowAudioBitRate      = "32000"
	config_VeryNarrowAudioBitRate  = "32000"
	config_WideResolution          = "720p" // 720p | 1080p
	config_NarrowResolution        = "720p"
	config_VeryNarrowResolution    = "720p"
	config_WideSpare2              = "sp2-a"
	config_NarrowSpare2            = "sp2-a"
	config_VeryNarrowSpare2        = "sp2-a"
	config_WideGain                = "-7"
	config_NarrowGain              = "-15"
	config_VeryNarrowGain          = "-20"
)

type (
	TxData_t struct {
		CurBand          string
		CurSymbolRate    string
		CurFrequency     string
		CurMode          string
		CurCodecs        string
		CurConstellation string
		CurFec           string
		CurVideoBitRate  string
		CurAudioBitRate  string
		CurResolution    string
		CurSpare2        string
		CurGain          string
		MarkerCentre     float32
		MarkerWidth      float32
		CurIsTuned       bool
		CurIsPtt         bool
	}

	selector_t struct {
		currIndex int
		lastIndex int
		list      []string
		value     string
	}
)

var (
	txData                TxData_t
	dataChan              chan TxData_t
	bandSelector          selector_t
	symbolRateSelector    selector_t
	frequencySelector     selector_t
	modeSelector          selector_t
	codecsSelector        selector_t
	constellationSelector selector_t
	fecSelector           selector_t
	videoBitRateSelector  selector_t
	audioBitRateSelector  selector_t
	resolutionSelector    selector_t
	spare2Selector        selector_t
	gainSelector          selector_t
	isTuned               bool
	isPtt                 bool
)

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
		"H264 G711u", "H264 ACC", "H265 G711u", "H265 ACC",
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
		"800", "810", "820", "830", "840", "850", "860", "870", "880", "890", "900", "910", "920", "930", "940", "950", "960", "970", "980", "990", "1000", "1100", "1200", "1300", "1400", "1500", "1600", "1700", "1800", "2000",
	}
	const_NARROW_VIDEO_BITRATE_LIST = []string{
		"250", "260", "270", "280", "290", "300", "310", "320", "330", "340", "350", "360", "370", "380", "390", "400", "410", "420", "430", "440", "450", "460", "470", "480", "490", "500",
	}
	const_VERY_NARROW_VIDEO_BITRATE_LIST = []string{
		"180", "190", "200", "210", "220", "230", "240", "250", "260", "270", "280", "290", "300",
	}
	const_WIDE_AUDIO_BITRATE_LIST = []string{
		"32000", "48000", "64000",
		// "32000", "48000", "64000", // 48000 unsupported
	}
	const_NARROW_AUDIO_BITRATE_LIST = []string{
		"32000", "48000", "64000",
		// "32000", "48000", "64000", // 48000 unsupported
	}
	const_VERY_NARROW_AUDIO_BITRATE_LIST = []string{
		"32000", "48000", "64000",
		// "32000", "48000", "64000", // 48000 unsupported
	}
	const_WIDE_RESOLUTION_LIST = []string{
		"720p", "1080p",
	}
	const_NARROW_RESOLUTION_LIST = []string{
		"720p", "1080p",
	}
	const_VERY_NARROW_RESOLUTION_LIST = []string{
		"720p", "1080p",
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
		"-16", "-15", "-14", "-13", "-12", "-11", "-10", "-9", "-8", "-7", "-6", "-5", "-4", "-3", "-2", "-1", "0",
	}
	const_NARROW_GAIN_LIST = []string{
		"-23", "-22", "-21", "-20", "-19", "-18", "-17", "-16", "-15", "-14", "-13", "-12", "-11", "-10", // "-9","-8","-7","-6","-5","-4","-3","-2","-1","0",
	}
	const_VERY_NARROW_GAIN_LIST = []string{
		"-23", "-22", "-21", "-20", "-19", "-18", "-17", "-16", "-15", "-14", "-13", "-12", "-11", "-10", //"-9","-8","-7","-6","-5","-4","-3","-2","-1","0",
	}

	wideSymbolRate          selector_t
	narrowSymbolRate        selector_t
	veryNarrowSymbolRate    selector_t
	wideFrequency           selector_t
	narrowFrequency         selector_t
	veryNarrowFrequency     selector_t
	wideMode                selector_t
	narrowMode              selector_t
	veryNarrowMode          selector_t
	wideCodecs              selector_t
	narrowCodecs            selector_t
	veryNarrowCodecs        selector_t
	wideConstellation       selector_t
	narrowConstellation     selector_t
	veryNarrowConstellation selector_t
	wideFec                 selector_t
	narrowFec               selector_t
	veryNarrowFec           selector_t
	wideVideoBitRate        selector_t
	narrowVideoBitRate      selector_t
	veryNarrowVideoBitRate  selector_t
	wideAudioBitRate        selector_t
	narrowAudioBitRate      selector_t
	veryNarrowAudioBitRate  selector_t
	wideResolution          selector_t
	narrowResolution        selector_t
	veryNarrowResolution    selector_t
	wideSpare2              selector_t
	narrowSpare2            selector_t
	veryNarrowSpare2        selector_t
	wideGain                selector_t
	narrowGain              selector_t
	veryNarrowGain          selector_t
)

type TxCmd_t int

const (
	CmdDecBand = iota
	CmdIncBand
	CmdDecSymbolRate
	CmdIncSymbolRate
	CmdDecFrequency
	CmdIncFrequency
	CmdDecMode
	CmdIncMode
	CmdDecCodecs
	CmdIncCodecs
	CmdDecConstellation
	CmdIncConstaellation
	CmdDecFec
	CmdIncFec
	CmdDecVideoBitRate
	CmdIncVideoBitRate
	CmdDecAudioBitRate
	CmdIncAudioBitRate
	CmdDecResolution
	CmdIncResolution
	CmdDecSpare2
	CmdIncSpare2
	CmdDecGain
	CmdIncGain
	CmdTune
	CmdPtt
)

func HandleCommands(ctx context.Context, cmdCh chan TxCmd_t, dataCh chan TxData_t) {

	dataChan = dataCh

	bandSelector = newSelector(const_BAND_LIST, config_Band)

	wideSymbolRate = newSelector(const_WIDE_SYMBOLRATE_LIST, config_WideSymbolrate)
	wideFrequency = newSelector(const_WIDE_FREQUENCY_LIST, config_WideFrequency)

	narrowSymbolRate = newSelector(const_NARROW_SYMBOLRATE_LIST, config_NarrowSymbolrate)
	narrowFrequency = newSelector(const_NARROW_FREQUENCY_LIST, config_NarrowFrequency)

	veryNarrowSymbolRate = newSelector(const_VERY_NARROW_SYMBOLRATE_LIST, config_VeryNarrowSymbolRate)
	veryNarrowFrequency = newSelector(const_VERY_NARROW_FREQUENCY_LIST, config_VeryNarrowFrequency)

	wideMode = newSelector(const_WIDE_MODE_LIST, config_WideMode)
	narrowMode = newSelector(const_NARROW_MODE_LIST, config_WideMode)
	veryNarrowMode = newSelector(const_VERY_NARROW_MODE_LIST, config_WideMode)

	wideCodecs = newSelector(const_WIDE_CODECS_LIST, config_WideCodecs)
	narrowCodecs = newSelector(const_NARROW_CODECS_LIST, config_NarrowCdecs)
	veryNarrowCodecs = newSelector(const_VERY_NARROW_CODECS_LIST, config_VeryNarrowCodecs)

	wideConstellation = newSelector(const_WIDE_CONSTELLATION_LIST, config_WideConstellation)
	narrowConstellation = newSelector(const_NARROW_CONSTELLATION_LIST, config_NarrowConstellation)
	veryNarrowConstellation = newSelector(const_VERY_NARROW_CONSTELLATION_LIST, config_VeryNarrorConstellation)

	wideFec = newSelector(const_WIDE_FEC_LIST, config_WideFec)
	narrowFec = newSelector(const_NARROW_FEC_LIST, config_NarrowFec)
	veryNarrowFec = newSelector(const_VERY_NARROW_FEC_LIST, config_VeryNarrowFec)

	wideVideoBitRate = newSelector(const_WIDE_VIDEO_BITRATE_LIST, config_WideVideoBitRate)
	narrowVideoBitRate = newSelector(const_NARROW_VIDEO_BITRATE_LIST, config_NarrowVideoBitRate)
	veryNarrowVideoBitRate = newSelector(const_VERY_NARROW_VIDEO_BITRATE_LIST, config_VeryNarrowVideoBitRate)

	wideAudioBitRate = newSelector(const_WIDE_AUDIO_BITRATE_LIST, config_WideAudioBitRate)
	narrowAudioBitRate = newSelector(const_NARROW_AUDIO_BITRATE_LIST, config_NarrowAudioBitRate)
	veryNarrowAudioBitRate = newSelector(const_VERY_NARROW_AUDIO_BITRATE_LIST, config_VeryNarrowAudioBitRate)

	wideResolution = newSelector(const_WIDE_RESOLUTION_LIST, config_WideResolution)
	narrowResolution = newSelector(const_NARROW_RESOLUTION_LIST, config_NarrowResolution)
	veryNarrowResolution = newSelector(const_VERY_NARROW_RESOLUTION_LIST, config_VeryNarrowResolution)

	wideSpare2 = newSelector(const_WIDE_SPARE2_LIST, config_WideSpare2)
	narrowSpare2 = newSelector(const_NARROW_SPARE2_LIST, config_NarrowSpare2)
	veryNarrowSpare2 = newSelector(const_VERY_NARROW_SPARE2_LIST, config_VeryNarrowSpare2)

	wideGain = newSelector(const_WIDE_GAIN_LIST, config_WideGain)
	narrowGain = newSelector(const_NARROW_GAIN_LIST, config_NarrowGain)
	veryNarrowGain = newSelector(const_VERY_NARROW_GAIN_LIST, config_VeryNarrowGain)

	switchBand()

	for {
		select {
		case <-ctx.Done():
			isPtt = pttSwitch.SetPtt(false)
			log.Printf("CANCEL ----- txControl has cancelled")
			return
		case command := <-cmdCh:
			switch command {
			case CmdDecBand:
				decBandSelector(&bandSelector)
			case CmdIncBand:
				incBandSelector(&bandSelector)
			case CmdDecSymbolRate:
				decSelector(&symbolRateSelector)
			case CmdIncSymbolRate:
				incSelector(&symbolRateSelector)
			case CmdDecFrequency:
				decSelector(&frequencySelector)
			case CmdIncFrequency:
				incSelector(&frequencySelector)
			case CmdDecMode:
				decSelector(&modeSelector)
			case CmdIncMode:
				incSelector(&modeSelector)
			case CmdDecCodecs:
				decSelector(&codecsSelector)
			case CmdIncCodecs:
				incSelector(&codecsSelector)
			case CmdDecConstellation:
				decSelector(&constellationSelector)
			case CmdIncConstaellation:
				incSelector(&constellationSelector)
			case CmdDecFec:
				decSelector(&fecSelector)
			case CmdIncFec:
				incSelector(&fecSelector)
			case CmdDecVideoBitRate:
				decSelector(&videoBitRateSelector)
			case CmdIncVideoBitRate:
				incSelector(&videoBitRateSelector)
			case CmdDecAudioBitRate:
				decSelector(&audioBitRateSelector)
			case CmdIncAudioBitRate:
				incSelector(&audioBitRateSelector)
			case CmdDecResolution:
				decSelector(&resolutionSelector)
			case CmdIncResolution:
				incSelector(&resolutionSelector)
			case CmdDecSpare2:
				decSelector(&spare2Selector)
			case CmdIncSpare2:
				incSelector(&spare2Selector)
			case CmdDecGain:
				decSelector(&gainSelector)
			case CmdIncGain:
				incSelector(&gainSelector)
			case CmdTune:
				setEncoder()
			case CmdPtt:
				setPtt()
			}
		}
	}
}

// called from the TUNE button
func setEncoder() {
	if !isTuned {
		plutoParam := plutoClient.PlConfig_t{
			Frequency:     frequencySelector.value,
			Mode:          modeSelector.value,
			Constellation: constellationSelector.value,
			SymbolRate:    symbolRateSelector.value,
			Fec:           fecSelector.value,
			Gain:          gainSelector.value,
		}
		plutoClient.SetParams(&plutoParam)

		encoderArgs := encoderClient.EncConfig_t{
			Codecs:       codecsSelector.value,
			AudioBitRate: audioBitRateSelector.value,
			VideoBitRate: videoBitRateSelector.value,
			Resolution:   resolutionSelector.value,
		}
		if err := encoderClient.SetParams(&encoderArgs); err != nil {
			log.Printf("ERROR TUNE: %s", err)
		}

		isTuned = true
	} else {
		isPtt = pttSwitch.SetPtt(false)
		isTuned = false
	}
	txData.CurIsTuned = isTuned
	txData.CurIsPtt = isPtt
	dataChan <- txData
}

// called from the PTT button
func setPtt() {
	if !isTuned {
		return
	}
	isPtt = pttSwitch.SetPtt(!isPtt)
	txData.CurIsPtt = isPtt
	dataChan <- txData
}

func incBandSelector(st *selector_t) {
	if st.currIndex < st.lastIndex {
		st.currIndex++
		st.value = st.list[st.currIndex]
		switchBand()
	}
}

func decBandSelector(st *selector_t) {
	if st.currIndex > 0 {
		st.currIndex--
		st.value = st.list[st.currIndex]
		switchBand()
	}
}

func incSelector(st *selector_t) {
	if st.currIndex < st.lastIndex {
		st.currIndex++
		st.value = st.list[st.currIndex]
		somethingChanged()
	}
}

func decSelector(st *selector_t) {
	if st.currIndex > 0 {
		st.currIndex--
		st.value = st.list[st.currIndex]
		somethingChanged()
	}
}

// TODO: add error chatching
func indexInList(list []string, with string) int {
	for i := range list {
		if list[i] == with {
			return i
		}
	}
	return 0
}

func newSelector(values []string, with string) selector_t {
	index := indexInList(values, with)
	st := selector_t{
		currIndex: index,
		lastIndex: len(values) - 1,
		list:      values,
		value:     values[index],
	}
	return st
}

func switchBand() { // TODO: should switch back to previosly use settings
	switch bandSelector.value {

	case const_BAND_LIST[0]: // wide
		symbolRateSelector = wideSymbolRate
		frequencySelector = wideFrequency
		modeSelector = wideMode
		codecsSelector = wideCodecs
		constellationSelector = wideConstellation
		fecSelector = wideFec
		videoBitRateSelector = wideVideoBitRate
		audioBitRateSelector = wideAudioBitRate
		resolutionSelector = wideResolution
		spare2Selector = wideSpare2
		gainSelector = wideGain
	case const_BAND_LIST[1]: // narrow
		symbolRateSelector = narrowSymbolRate
		frequencySelector = narrowFrequency
		modeSelector = narrowMode
		codecsSelector = narrowCodecs
		constellationSelector = narrowConstellation
		fecSelector = narrowFec
		videoBitRateSelector = narrowVideoBitRate
		audioBitRateSelector = narrowAudioBitRate
		resolutionSelector = narrowResolution
		spare2Selector = narrowSpare2
		gainSelector = narrowGain
	case const_BAND_LIST[2]: // very narrow
		symbolRateSelector = veryNarrowSymbolRate
		frequencySelector = veryNarrowFrequency
		modeSelector = veryNarrowMode
		codecsSelector = veryNarrowCodecs
		constellationSelector = veryNarrowConstellation
		fecSelector = veryNarrowFec
		videoBitRateSelector = veryNarrowVideoBitRate
		audioBitRateSelector = veryNarrowAudioBitRate
		resolutionSelector = veryNarrowResolution
		spare2Selector = veryNarrowSpare2
		gainSelector = veryNarrowGain
	}
	somethingChanged()
}

func somethingChanged() {
	pttSwitch.SetPtt(false)
	isPtt = false
	isTuned = false
	txData.CurBand = bandSelector.value
	txData.CurSymbolRate = symbolRateSelector.value
	txData.CurFrequency = frequencySelector.value
	txData.CurMode = modeSelector.value
	txData.CurCodecs = codecsSelector.value
	txData.CurConstellation = constellationSelector.value
	txData.CurFec = fecSelector.value
	txData.CurVideoBitRate = videoBitRateSelector.value
	txData.CurAudioBitRate = audioBitRateSelector.value
	txData.CurResolution = resolutionSelector.value
	txData.CurSpare2 = spare2Selector.value
	txData.CurGain = gainSelector.value
	txData.MarkerCentre = const_frequencyCentre[frequencySelector.value] / 9.18 // NOTE: 9.18 is a temporary kludge
	txData.MarkerWidth = const_symbolRateWidth[symbolRateSelector.value]
	txData.CurIsTuned = isTuned
	txData.CurIsPtt = isPtt
	dataChan <- txData
}
