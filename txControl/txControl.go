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

type (
	TuConfig_t struct {
		Band                    string
		WideFrequency           string
		WideSymbolrate          string
		NarrowFrequency         string
		NarrowSymbolrate        string
		VeryNarrowFrequency     string
		VeryNarrowSymbolRate    string
		WideMode                string
		NarrowMode              string
		VeryNarrowMode          string
		WideCodecs              string
		NarrowCdecs             string
		VeryNarrowCodecs        string
		WideConstellation       string
		NarrowConstellation     string
		VeryNarrorConstellation string
		WideFec                 string
		NarrowFec               string
		VeryNarrowFec           string
		WideVideoBitRate        string
		NarrowVideoBitRate      string
		VeryNarrowVideoBitRate  string
		WideAudioBitRate        string
		NarrowAudioBitRate      string
		VeryNarrowAudioBitRate  string
		WideResolution          string
		NarrowResolution        string
		VeryNarrowResolution    string
		WideSpare2              string
		NarrowSpare2            string
		VeryNarrowSpare2        string
		WideGain                string
		NarrowGain              string
		VeryNarrowGain          string
	}
	Selector_t struct {
		currIndex int
		lastIndex int
		list      []string
		Value     string
	}
	TuData_t struct {
		MarkerCentre float32
		MarkerWidth  float32
	}
)

var (
	tuCmd TuCmd_t
	// cmdChan chan TuCmd_t

	tuData   TuData_t
	dataChan chan TuData_t

	Band          Selector_t
	SymbolRate    Selector_t
	Frequency     Selector_t
	Mode          Selector_t
	Codecs        Selector_t
	Constellation Selector_t
	Fec           Selector_t
	VideoBitRate  Selector_t
	AudioBitRate  Selector_t
	Resolution    Selector_t
	Spare2        Selector_t
	Gain          Selector_t
	IsTuned       bool
	IsPtt         bool
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
		"380", "390", "400", "410", "420", "430", "440", "450", "460", "470", "480", "490", "410", "420", "430", "440", "450", "460", "470", "480", "490", "500",
	}
	const_NARROW_VIDEO_BITRATE_LIST = []string{
		"180", "190", "200", "210", "220", "230", "240", "250", "260", "270", "280", "290", "300", "310", "320", "330", "340", "350", "360", "370", "380", "390",
	}
	const_VERY_NARROW_VIDEO_BITRATE_LIST = []string{
		"180", "190", "200", "210", "220", "230", "240", "250", "260", "270", "280", "290", "300", "310", "320", "330", "340", "350", "360", "370", "380", "390",
	}
	const_WIDE_AUDIO_BITRATE_LIST = []string{
		"32000", "64000",
		// "32000", "48000", "64000", // 48000 unsupported
	}
	const_NARROW_AUDIO_BITRATE_LIST = []string{
		"32000", "64000",
		// "32000", "48000", "64000", // 48000 unsupported
	}
	const_VERY_NARROW_AUDIO_BITRATE_LIST = []string{
		"32000", "64000",
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
		"-23", "-22", "-21", "-20", "-19", "-18", "-17", "-16", "-15", "-14", "-13", "-12", "-11", "-10", // "-9","-8","-7","-6","-5","-4","-3","-2","-1","0",
	}
	const_NARROW_GAIN_LIST = []string{
		"-23", "-22", "-21", "-20", "-19", "-18", "-17", "-16", "-15", "-14", "-13", "-12", "-11", "-10", // "-9","-8","-7","-6","-5","-4","-3","-2","-1","0",
	}
	const_VERY_NARROW_GAIN_LIST = []string{
		"-23", "-22", "-21", "-20", "-19", "-18", "-17", "-16", "-15", "-14", "-13", "-12", "-11", "-10", //"-9","-8","-7","-6","-5","-4","-3","-2","-1","0",
	}

	wideSymbolRate          Selector_t
	narrowSymbolRate        Selector_t
	veryNarrowSymbolRate    Selector_t
	wideFrequency           Selector_t
	narrowFrequency         Selector_t
	veryNarrowFrequency     Selector_t
	wideMode                Selector_t
	narrowMode              Selector_t
	veryNarrowMode          Selector_t
	wideCodecs              Selector_t
	narrowCodecs            Selector_t
	veryNarrowCodecs        Selector_t
	wideConstellation       Selector_t
	narrowConstellation     Selector_t
	veryNarrowConstellation Selector_t
	wideFec                 Selector_t
	narrowFec               Selector_t
	veryNarrowFec           Selector_t
	wideVideoBitRate        Selector_t
	narrowVideoBitRate      Selector_t
	veryNarrowVideoBitRate  Selector_t
	wideAudioBitRate        Selector_t
	narrowAudioBitRate      Selector_t
	veryNarrowAudioBitRate  Selector_t
	wideResolution          Selector_t
	narrowResolution        Selector_t
	veryNarrowResolution    Selector_t
	wideSpare2              Selector_t
	narrowSpare2            Selector_t
	veryNarrowSpare2        Selector_t
	wideGain                Selector_t
	narrowGain              Selector_t
	veryNarrowGain          Selector_t
)

type TuCmd_t int

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

func HandleCommands(ctx context.Context, cfg TuConfig_t, cmdCh chan TuCmd_t, dataCh chan TuData_t) {
	// cmdChan = cmdCh
	dataChan = dataCh

	Band = newSelector(const_BAND_LIST, cfg.Band)

	wideSymbolRate = newSelector(const_WIDE_SYMBOLRATE_LIST, cfg.WideSymbolrate)
	wideFrequency = newSelector(const_WIDE_FREQUENCY_LIST, cfg.WideFrequency)

	narrowSymbolRate = newSelector(const_NARROW_SYMBOLRATE_LIST, cfg.NarrowSymbolrate)
	narrowFrequency = newSelector(const_NARROW_FREQUENCY_LIST, cfg.NarrowFrequency)

	veryNarrowSymbolRate = newSelector(const_VERY_NARROW_SYMBOLRATE_LIST, cfg.VeryNarrowSymbolRate)
	veryNarrowFrequency = newSelector(const_VERY_NARROW_FREQUENCY_LIST, cfg.VeryNarrowFrequency)

	wideMode = newSelector(const_WIDE_MODE_LIST, cfg.WideMode)
	narrowMode = newSelector(const_NARROW_MODE_LIST, cfg.WideMode)
	veryNarrowMode = newSelector(const_VERY_NARROW_MODE_LIST, cfg.WideMode)

	wideCodecs = newSelector(const_WIDE_CODECS_LIST, cfg.WideCodecs)
	narrowCodecs = newSelector(const_NARROW_CODECS_LIST, cfg.NarrowCdecs)
	veryNarrowCodecs = newSelector(const_VERY_NARROW_CODECS_LIST, cfg.VeryNarrowCodecs)

	wideConstellation = newSelector(const_WIDE_CONSTELLATION_LIST, cfg.WideConstellation)
	narrowConstellation = newSelector(const_NARROW_CONSTELLATION_LIST, cfg.NarrowConstellation)
	veryNarrowConstellation = newSelector(const_VERY_NARROW_CONSTELLATION_LIST, cfg.VeryNarrorConstellation)

	wideFec = newSelector(const_WIDE_FEC_LIST, cfg.WideFec)
	narrowFec = newSelector(const_NARROW_FEC_LIST, cfg.NarrowFec)
	veryNarrowFec = newSelector(const_VERY_NARROW_FEC_LIST, cfg.VeryNarrowFec)

	wideVideoBitRate = newSelector(const_WIDE_VIDEO_BITRATE_LIST, cfg.WideVideoBitRate)
	narrowVideoBitRate = newSelector(const_NARROW_VIDEO_BITRATE_LIST, cfg.NarrowVideoBitRate)
	veryNarrowVideoBitRate = newSelector(const_VERY_NARROW_VIDEO_BITRATE_LIST, cfg.VeryNarrowVideoBitRate)

	wideAudioBitRate = newSelector(const_WIDE_AUDIO_BITRATE_LIST, cfg.WideAudioBitRate)
	narrowAudioBitRate = newSelector(const_NARROW_AUDIO_BITRATE_LIST, cfg.NarrowAudioBitRate)
	veryNarrowAudioBitRate = newSelector(const_VERY_NARROW_AUDIO_BITRATE_LIST, cfg.VeryNarrowAudioBitRate)

	wideResolution = newSelector(const_WIDE_RESOLUTION_LIST, cfg.WideResolution)
	narrowResolution = newSelector(const_NARROW_RESOLUTION_LIST, cfg.NarrowResolution)
	veryNarrowResolution = newSelector(const_VERY_NARROW_RESOLUTION_LIST, cfg.VeryNarrowResolution)

	wideSpare2 = newSelector(const_WIDE_SPARE2_LIST, cfg.WideSpare2)
	narrowSpare2 = newSelector(const_NARROW_SPARE2_LIST, cfg.NarrowSpare2)
	veryNarrowSpare2 = newSelector(const_VERY_NARROW_SPARE2_LIST, cfg.VeryNarrowSpare2)

	wideGain = newSelector(const_WIDE_GAIN_LIST, cfg.WideGain)
	narrowGain = newSelector(const_NARROW_GAIN_LIST, cfg.NarrowGain)
	veryNarrowGain = newSelector(const_VERY_NARROW_GAIN_LIST, cfg.VeryNarrowGain)

	switchBand()

	for {
		select {
		case <-ctx.Done():
			IsPtt = pttSwitch.SetPtt(false)
			log.Printf("INFO ----- txControl has stopped")
			return
		case tuCmd = <-cmdCh:
			switch tuCmd {
			case CmdDecBand:
				DecBandSelector(&Band)
			case CmdIncBand:
				IncBandSelector(&Band)
			case CmdDecSymbolRate:
				DecSelector(&SymbolRate)
			case CmdIncSymbolRate:
				IncSelector(&SymbolRate)
			case CmdDecFrequency:
				DecSelector(&Frequency)
			case CmdIncFrequency:
				IncSelector(&Frequency)
			case CmdDecMode:
				DecSelector(&Mode)
			case CmdIncMode:
				IncSelector(&Mode)
			case CmdDecCodecs:
				DecSelector(&Codecs)
			case CmdIncCodecs:
				IncSelector(&Codecs)
			case CmdDecConstellation:
				DecSelector(&Constellation)
			case CmdIncConstaellation:
				IncSelector(&Constellation)
			case CmdDecFec:
				DecSelector(&Fec)
			case CmdIncFec:
				IncSelector(&Fec)
			case CmdDecVideoBitRate:
				DecSelector(&VideoBitRate)
			case CmdIncVideoBitRate:
				IncSelector(&VideoBitRate)
			case CmdDecAudioBitRate:
				DecSelector(&AudioBitRate)
			case CmdIncAudioBitRate:
				IncSelector(&AudioBitRate)
			case CmdDecResolution:
				DecSelector(&Resolution)
			case CmdIncResolution:
				IncSelector(&Resolution)
			case CmdDecSpare2:
				DecSelector(&Spare2)
			case CmdIncSpare2:
				IncSelector(&Spare2)
			case CmdDecGain:
				DecSelector(&Gain)
			case CmdIncGain:
				IncSelector(&Gain)
			case CmdTune:
				Tune()
			case CmdPtt:
				Ptt()
			}
		}
		log.Printf("...")
	}
}

// func Stop() {
// 	log.Printf("INFO Tuner will stop... - NOT IMPLEMENTED")
// 	IsPtt = pttSwitch.SetPtt(false)
// 	log.Printf("INFO Tuner has stopped")
// }

// API
func Tune() {
	if !IsTuned {
		plutoParam := plutoClient.PlConfig_t{
			Frequency:     Frequency.Value,
			Mode:          Mode.Value,
			Constellation: Constellation.Value,
			SymbolRate:    SymbolRate.Value,
			Fec:           Fec.Value,
			Gain:          Gain.Value,
		}
		plutoClient.SetParams(&plutoParam)

		encoderArgs := encoderClient.EncConfig_t{
			Codecs:       Codecs.Value,
			AudioBitRate: AudioBitRate.Value,
			VideoBitRate: VideoBitRate.Value,
			Resolution:   Resolution.Value,
		}
		if err := encoderClient.SetParams(&encoderArgs); err != nil {
			log.Printf("ERROR TUNE: %s", err)
		}

		IsTuned = true
	} else {
		// if IsPtt {
		IsPtt = pttSwitch.SetPtt(false)
		// 	// IsPtt = false
		// }
		IsTuned = false
	}
	// log.Printf("INFO IsTuned is %v", IsTuned)
}

// API
func Ptt() {
	if !IsTuned {
		return
	}
	if IsPtt {
		pttSwitch.SetPtt(false)
		IsPtt = false
	} else {
		pttSwitch.SetPtt(true)
		IsPtt = true
	}
}

// API
func IncBandSelector(st *Selector_t) {
	if st.currIndex < st.lastIndex {
		st.currIndex++
		st.Value = st.list[st.currIndex]
		switchBand()
	}
}

// API
func DecBandSelector(st *Selector_t) {
	if st.currIndex > 0 {
		st.currIndex--
		st.Value = st.list[st.currIndex]
		switchBand()
	}
}

// API
func IncSelector(st *Selector_t) {
	if st.currIndex < st.lastIndex {
		st.currIndex++
		st.Value = st.list[st.currIndex]
		somethingChanged()
	}
}

// API
func DecSelector(st *Selector_t) {
	if st.currIndex > 0 {
		st.currIndex--
		st.Value = st.list[st.currIndex]
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

func newSelector(values []string, with string) Selector_t {
	index := indexInList(values, with)
	st := Selector_t{
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
		Resolution = wideResolution
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
		Resolution = narrowResolution
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
		Resolution = veryNarrowResolution
		Spare2 = veryNarrowSpare2
		Gain = veryNarrowGain
	}
	somethingChanged()
}

func somethingChanged() {
	pttSwitch.SetPtt(false)
	IsPtt = false
	IsTuned = false
	tuData.MarkerCentre = const_frequencyCentre[Frequency.Value] / 9.18 // NOTE: 9.18 is a temporary kludge
	tuData.MarkerWidth = const_symbolRateWidth[SymbolRate.Value]
	dataChan <- tuData
}
