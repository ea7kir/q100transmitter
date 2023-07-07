/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package tuner

import (
	"q100transmitter/hev10"
	"q100transmitter/logger"
)

type (
	TuConfig struct {
		Band                 string
		WideFrequency        string
		WideSymbolrate       string
		NarrowFrequency      string
		NarrowSymbolrate     string
		VeryNarrowFrequency  string
		VeryNarrowSymbolRate string
	}
)

var (
	Band       Selector
	SymbolRate Selector
	Frequency  Selector
	IsTuned    = false
	IsPtt      = false
)

func Intitialize(tuc TuConfig) {
	// find index in list
	Band = newSelector(const_BAND_LIST, indexInList(const_BAND_LIST, tuc.Band))
	beaconSymbolRate = newSelector(const_BEACON_SYMBOLRATE_LIST, 0)
	beaconFrequency = newSelector(const_BEACON_FREQUENCY_LIST, 0)
	wideSymbolRate = newSelector(const_WIDE_SYMBOLRATE_LIST, indexInList(const_WIDE_SYMBOLRATE_LIST, tuc.WideSymbolrate))
	wideFrequency = newSelector(const_WIDE_FREQUENCY_LIST, indexInList(const_WIDE_FREQUENCY_LIST, tuc.WideFrequency))
	narrowSymbolRate = newSelector(const_NARROW_SYMBOLRATE_LIST, indexInList(const_NARROW_SYMBOLRATE_LIST, tuc.NarrowSymbolrate))
	narrowFrequency = newSelector(const_NARROW_FREQUENCY_LIST, indexInList(const_NARROW_FREQUENCY_LIST, tuc.NarrowFrequency))
	veryNarrowSymbolRate = newSelector(const_VERY_NARROW_SYMBOLRATE_LIST, indexInList(const_VERY_NARROW_SYMBOLRATE_LIST, tuc.NarrowSymbolrate))
	veryNarrowFrequency = newSelector(const_VERY_NARROW_FREQUENCY_LIST, indexInList(const_VERY_NARROW_FREQUENCY_LIST, tuc.VeryNarrowFrequency))
}

func Start() {
	logger.Info.Printf("Tuner will start...")
	switchBand()
	logger.Info.Printf("Tuner has started")
}

func Stop() {
	logger.Info.Printf("Tuner will stop... - DOES NOTHING")
	//
	logger.Info.Printf("Tuner has stopped")
}

func Tune() {
	if IsTuned {
		hev10.UnConfig()
		IsTuned = false
	} else {
		hev10.Config()
		IsTuned = true
	}
	// logger.Info.Printf("IsTuned is %v", IsTuned)
}

func Ptt() {
	if IsPtt {
		IsPtt = false
	} else {
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
		if IsTuned {
			IsTuned = false
			logger.Info.Printf("IsTuned is %v", IsTuned)

		}
		st.currIndex++
		st.Value = st.list[st.currIndex]
		switchBand()
	}
}

func DecBandSelector(st *Selector) {
	if st.currIndex > 0 {
		st.currIndex--
		st.Value = st.list[st.currIndex]
		switchBand()
	}
}

func IncSelector(st *Selector) {
	if st.currIndex < st.lastIndex {
		st.currIndex++
		st.Value = st.list[st.currIndex]
	}
}

func DecSelector(st *Selector) {
	if st.currIndex > 0 {
		st.currIndex--
		st.Value = st.list[st.currIndex]
	}
}

func IncFrequencySelector(st *Selector) {
	if st.currIndex < st.lastIndex {
		st.currIndex++
		st.Value = st.list[st.currIndex]
	}
}

func DecFrequencySelector(st *Selector) {
	if st.currIndex > 0 {
		st.currIndex--
		st.Value = st.list[st.currIndex]
	}
}
