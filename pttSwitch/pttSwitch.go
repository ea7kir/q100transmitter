/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package pttSwitch

import (
	"q100transmitter/logger"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

const (
	// Analog Devices HVC349 RF Switch
	// using RFC as the input, RF2 as the output
	CTRL_PIN = rpi.J8p18 // pin 18 GPIO_24
	EN_PIN   = rpi.J8p16 // pin 16 GPIO_23 (GND pin 14, pin 17 3.3v)
	LOW      = 0
	HIGH     = 1
)

var (
	hvc349Control *gpiod.Line
	hvc349Enable  *gpiod.Line
)

// API
func Initialize() {
	hvc349ControlLine, err := gpiod.RequestLine("gpiochip0", CTRL_PIN, gpiod.AsOutput(0))
	if err != nil {
		panic(err)
	}
	hvc349Control = hvc349ControlLine
	hvc349ControlLine.SetValue(HIGH)
	hvc349EnableLine, err := gpiod.RequestLine("gpiochip0", EN_PIN, gpiod.AsOutput(0))
	if err != nil {
		panic(err)
	}
	hvc349Enable = hvc349EnableLine
	hvc349Enable.SetValue(LOW)
}

// API
func Stop() {
	SetPtt(false)
	hvc349Control.SetValue(HIGH)
	hvc349Enable.SetValue(HIGH)

	hvc349Control.Reconfigure(gpiod.AsInput)
	hvc349Control.Close()
	hvc349Enable.Reconfigure(gpiod.AsInput)
	hvc349Enable.Close()
}

// API
func SetPtt(tx bool) bool {
	switch tx {
	case true:
		hvc349Control.SetValue(LOW)
		logger.Info("PTT is %v", "Enabled")
	case false:
		hvc349Control.SetValue(HIGH)
		logger.Info("PTT is %v", "Disabled")
	}
	return tx
}

// there is a command app called: raspi-gpio

/* Analog Devices HMC349 RF Switch

# Analog Devices HVC349 RF Switch
# using RFC as the input, RF2 as the output
EN_PIN               = 23 # pin 16 GPIO_23 (GND pin 14, pin 17 3.3v)
CTRL_PIN             = 24 # pin 18 GPIO_24
HIGH                  = 1
LOW                   = 0

_rf_pi = None

# HMC349 using RFC as the input, RF2 as the output

def _hvc349_HMC349(gpio, state):
    _rf_pi.write(gpio, state)

def _config_HMC349(en_gpio, ctrl_gpio):
    _rf_pi.set_mode(ctrl_gpio, pigpio.OUTPUT)
    _rf_pi.write(ctrl_gpio, HIGH)
    _rf_pi.set_mode(en_gpio, pigpio.OUTPUT)
    _rf_pi.write(en_gpio, LOW)

def configure_rf_hvc349es(pi):
    global _rf_pi
    _rf_pi = pi
    _config_HMC349(EN_PIN, CTRL_PIN)

def shutdown_rf_hvc349es():
    _hvc349_HMC349(CTRL_PIN, HIGH)
    _hvc349_HMC349(EN_PIN, HIGH)

def hvc349_rf_hvc349_On():
    #print("SWITCHING ON RF SWITCH", flush=True)
    _hvc349_HMC349(CTRL_PIN, LOW)

def hvc349_rf_hvc349_Off():
    #print("SWITCHING OFF RF SWITCH", flush=True)
    _hvc349_HMC349(CTRL_PIN, HIGH)

*/
