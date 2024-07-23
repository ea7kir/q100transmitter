/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package pttSwitch

import (
	"log"

	"github.com/warthog618/go-gpiocdev"
	"github.com/warthog618/go-gpiocdev/device/rpi"
)

const (
	// Analog Devices HVC349 RF Switch
	// using RFC as the input, RF2 as the output, RF1 with 50ohm dummy load

	// GND_PIN to RPi pin 14 Ground            - Brown
	// VCC_PIN to RPi pin 17 3V3               - Red

	kCTRL_PIN = rpi.J8p18 // RPi pin 18 GPIO_24 - Yellow

	// Remote mute to the LU Meter project
	// using the TX-Remote PCB

	// GND_PIN to RPi pin 39 Ground			   - Brown

	kMUTE_PIN = rpi.J8p40 // RPi pin 16 GPIO_21 - Orenage

	kLOW  = 0
	kHIGH = 1
)

var (
	hvc349Control *gpiocdev.Line
	muteEnable    *gpiocdev.Line
)

// API
func Start() {
	hvc349ControlLine, err := gpiocdev.RequestLine("gpiochip0", kCTRL_PIN, gpiocdev.AsOutput(0))
	if err != nil {
		panic(err)
	}
	hvc349Control = hvc349ControlLine
	hvc349ControlLine.SetValue(kHIGH)
	muteEnableLine, err := gpiocdev.RequestLine("gpiochip0", kMUTE_PIN, gpiocdev.AsOutput(0))
	if err != nil {
		panic(err)
	}
	muteEnable = muteEnableLine
	muteEnable.SetValue(kLOW)
	log.Printf("pttSwitch has started")
}

// API
func Stop() {
	SetPtt(false)
	hvc349Control.SetValue(kHIGH)
	muteEnable.SetValue(kLOW)

	hvc349Control.Reconfigure(gpiocdev.AsInput)
	hvc349Control.Close()
	muteEnable.Reconfigure(gpiocdev.AsInput)
	muteEnable.Close()
	log.Printf("pttSwitch has stopped")
}

// API
func SetPtt(tx bool) bool {
	switch tx {
	case true:
		hvc349Control.SetValue(kLOW)
		muteEnable.SetValue(kHIGH)
		log.Printf("INFO PTT is %v", "Enabled")
	case false:
		hvc349Control.SetValue(kHIGH)
		muteEnable.SetValue(kLOW)
		// log.Printf("INFO PTT is %v", "Disabled") // too much logging, because called on when any button is pressed
	}
	return tx
}
