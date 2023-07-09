/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package tuner

import (
	"q100transmitter/hev10"
	"q100transmitter/logger"
	"q100transmitter/pluto"
	"q100transmitter/rfSwitch"
)

type (
	TuConfig struct {
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
		WideSpare1              string
		NarrowSpare1            string
		VeryNarrowSpare1        string
		WideSpare2              string
		NarrowSpare2            string
		VeryNarrowSpare2        string
		WideGain                string
		NarrowGain              string
		VeryNarrowGain          string
	}
)

var (
	Band          Selector
	SymbolRate    Selector
	Frequency     Selector
	Mode          Selector
	Codecs        Selector
	Constellation Selector
	Fec           Selector
	VideoBitRate  Selector
	AudioBitRate  Selector
	Spare1        Selector
	Spare2        Selector
	Gain          Selector

	IsTuned = false
	IsPtt   = false
)

func Intitialize(tuc TuConfig) {
	// find index in list
	Band = newSelector(const_BAND_LIST, tuc.Band)
	wideSymbolRate = newSelector(const_WIDE_SYMBOLRATE_LIST, tuc.WideSymbolrate)
	wideFrequency = newSelector(const_WIDE_FREQUENCY_LIST, tuc.WideFrequency)
	narrowSymbolRate = newSelector(const_NARROW_SYMBOLRATE_LIST, tuc.NarrowSymbolrate)
	narrowFrequency = newSelector(const_NARROW_FREQUENCY_LIST, tuc.NarrowFrequency)
	veryNarrowSymbolRate = newSelector(const_VERY_NARROW_SYMBOLRATE_LIST, tuc.VeryNarrowSymbolRate)
	veryNarrowFrequency = newSelector(const_VERY_NARROW_FREQUENCY_LIST, tuc.VeryNarrowFrequency)

	wideMode = newSelector(const_WIDE_MODE_LIST, tuc.WideMode)
	narrowMode = newSelector(const_NARROW_MODE_LIST, tuc.WideMode)
	veryNarrowMode = newSelector(const_VERY_NARROW_MODE_LIST, tuc.WideMode)

	wideCodecs = newSelector(const_WIDE_CODECS_LIST, tuc.WideCodecs)
	narrowCodecs = newSelector(const_NARROW_CODECS_LIST, tuc.NarrowCdecs)
	veryNarrowCodecs = newSelector(const_VERY_NARROW_CODECS_LIST, tuc.VeryNarrowCodecs)

	wideConstellation = newSelector(const_WIDE_CONSTELLATION_LIST, tuc.WideConstellation)
	narrowConstellation = newSelector(const_NARROW_CONSTELLATION_LIST, tuc.NarrowConstellation)
	veryNarrowConstellation = newSelector(const_VERY_NARROW_CONSTELLATION_LIST, tuc.VeryNarrorConstellation)

	wideFec = newSelector(const_WIDE_FEC_LIST, tuc.WideFec)
	narrowFec = newSelector(const_NARROW_FEC_LIST, tuc.NarrowFec)
	veryNarrowFec = newSelector(const_VERY_NARROW_FEC_LIST, tuc.VeryNarrowFec)

	wideVideoBitRate = newSelector(const_WIDE_VIDEO_BITRATE_LIST, tuc.WideVideoBitRate)
	narrowVideoBitRate = newSelector(const_NARROW_VIDEO_BITRATE_LIST, tuc.NarrowVideoBitRate)
	veryNarrowVideoBitRate = newSelector(const_VERY_NARROW_VIDEO_BITRATE_LIST, tuc.VeryNarrowVideoBitRate)

	wideAudioBitRate = newSelector(const_WIDE_AUDIO_BITRATE_LIST, tuc.WideAudioBitRate)
	narrowAudioBitRate = newSelector(const_NARROW_AUDIO_BITRATE_LIST, tuc.NarrowAudioBitRate)
	veryNarrowAudioBitRate = newSelector(const_VERY_NARROW_AUDIO_BITRATE_LIST, tuc.VeryNarrowAudioBitRate)

	wideSpare1 = newSelector(const_WIDE_SPARE1_LIST, tuc.WideSpare1)
	narrowSpare1 = newSelector(const_NARROW_SPARE1_LIST, tuc.NarrowSpare1)
	veryNarrowSpare1 = newSelector(const_VERY_NARROW_SPARE1_LIST, tuc.VeryNarrowSpare1)

	wideSpare2 = newSelector(const_WIDE_SPARE2_LIST, tuc.WideSpare2)
	narrowSpare2 = newSelector(const_NARROW_SPARE2_LIST, tuc.NarrowSpare2)
	veryNarrowSpare2 = newSelector(const_VERY_NARROW_SPARE2_LIST, tuc.VeryNarrowSpare2)

	wideGain = newSelector(const_WIDE_GAIN_LIST, tuc.WideGain)
	narrowGain = newSelector(const_NARROW_GAIN_LIST, tuc.NarrowGain)
	veryNarrowGain = newSelector(const_VERY_NARROW_GAIN_LIST, tuc.VeryNarrowGain)

	switchBand()
}

// func Start() {
// 	// logger.Info.Printf("Tuner will start...")
// 	// switchBand()
// 	// logger.Info.Printf("Tuner has started")
// }

func Stop() {
	logger.Info.Printf("Tuner will stop...")
	IsPtt = rfSwitch.SetPtt(false)
	logger.Info.Printf("Tuner has stopped")
}

func Tune() {
	if IsTuned {
		// if IsPtt {
		IsPtt = rfSwitch.SetPtt(false)
		// 	// IsPtt = false
		// }
		IsTuned = false
	} else {
		hev10.SetParams(nil)
		pluto.SetParams(nil)
		IsTuned = true
	}
	// logger.Info.Printf("IsTuned is %v", IsTuned)
}

func Ptt() {
	if !IsTuned {
		return
	}
	if IsPtt {
		rfSwitch.SetPtt(false)
		IsPtt = false
	} else {
		rfSwitch.SetPtt(true)
		IsPtt = true
	}
}

type Selector struct {
	currIndex int
	lastIndex int
	list      []string
	Value     string
}

func IncBandSelector(st *Selector) {
	if st.currIndex < st.lastIndex {
		// if IsTuned {
		// 	IsTuned = false
		// 	logger.Info.Printf("IsTuned is %v", IsTuned)

		// }
		st.currIndex++
		st.Value = st.list[st.currIndex]
		switchBand()
		somethingChanged()
	}
}

func DecBandSelector(st *Selector) {
	if st.currIndex > 0 {
		st.currIndex--
		st.Value = st.list[st.currIndex]
		switchBand()
		somethingChanged()
	}
}

func IncSelector(st *Selector) {
	if st.currIndex < st.lastIndex {
		st.currIndex++
		st.Value = st.list[st.currIndex]
		somethingChanged()
	}
}

func DecSelector(st *Selector) {
	if st.currIndex > 0 {
		st.currIndex--
		st.Value = st.list[st.currIndex]
		somethingChanged()
	}
}
