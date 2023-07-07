/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package main

import (
	"os"
	"q100transmitter/spReader"
	"q100transmitter/tuner"
)

// application directory
// NOTE: this only works if we are are already in the correct folder
const appFolder = "/home/pi/Q100/q100transmitter-v1/"

// configuration data
var (
	spectrumConfig = spReader.SpConfig{
		Url: "wss://eshail.batc.org.uk/wb/fft/fft_ea7kirsatcontroller:443/",
	}

	tuConfig = tuner.TuConfig{
		Band:                 "Narrow",
		WideFrequency:        "10494.75 / 09",
		WideSymbolrate:       "1000",
		NarrowFrequency:      "10499.25 / 27",
		NarrowSymbolrate:     "333",
		VeryNarrowFrequency:  "10496.00 / 14",
		VeryNarrowSymbolRate: "125",
	}
)

// local data
var (
	spData    spReader.SpData
	spChannel = make(chan spReader.SpData, 5)
)

func main() {
	os.Setenv("DISPLAY", ":0") // required for X11

	spReader.Intitialize(spectrumConfig, spChannel)
	spReader.Start()

	tuner.Intitialize(tuConfig)
	tuner.Start()
}
