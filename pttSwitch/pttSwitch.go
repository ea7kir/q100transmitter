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
	RF_SWITCH_CTRL_GPIO = rpi.J8p18 // pin 18 GPIO_24
	RF_SWITCH_EN_GPIO   = rpi.J8p16 // pin 16 GPIO_23 (GND pin 14, pin 17 3.3v)
	RF_SWITCH_LOW       = 0
	RF_SWITCH_HIGH      = 1
)

var (
	switchControl *gpiod.Line
	switchEnable  *gpiod.Line
)

// API
func Initialize() {
	switchControlLine, err := gpiod.RequestLine("gpiochip0", RF_SWITCH_CTRL_GPIO, gpiod.AsOutput(0))
	if err != nil {
		panic(err)
	}
	switchControl = switchControlLine
	switchControlLine.SetValue(RF_SWITCH_HIGH)
	switchEnableLine, err := gpiod.RequestLine("gpiochip0", RF_SWITCH_EN_GPIO, gpiod.AsOutput(0))
	if err != nil {
		panic(err)
	}
	switchEnable = switchEnableLine
	switchEnable.SetValue(RF_SWITCH_LOW)
}

// API
func Stop() {
	SetPtt(false)
	switchControl.SetValue(RF_SWITCH_HIGH)
	switchEnable.SetValue(RF_SWITCH_HIGH)

	switchControl.Reconfigure(gpiod.AsInput)
	switchControl.Close()
	switchEnable.Reconfigure(gpiod.AsInput)
	switchEnable.Close()
}

// API
func SetPtt(tx bool) bool {
	switch tx {
	case true:
		switchControl.SetValue(RF_SWITCH_LOW)
		logger.Info.Printf("PTT is %v", "Enabled")
	case false:
		switchControl.SetValue(RF_SWITCH_HIGH)
		logger.Info.Printf("PTT is %v", "Disabled")
	}
	return tx
}

// there is a command app called: raspi-gpio

/* Analog Devices HMC349 RF Switch

# Analog Devices HVC349 RF Switch
# using RFC as the input, RF2 as the output
RF_SWITCH_EN_GPIO               = 23 # pin 16 GPIO_23 (GND pin 14, pin 17 3.3v)
RF_SWITCH_CTRL_GPIO             = 24 # pin 18 GPIO_24
RF_SWITCH_HIGH                  = 1
RF_SWITCH_LOW                   = 0

_rf_pi = None

# HMC349 using RFC as the input, RF2 as the output

def _switch_HMC349(gpio, state):
    _rf_pi.write(gpio, state)

def _config_HMC349(en_gpio, ctrl_gpio):
    _rf_pi.set_mode(ctrl_gpio, pigpio.OUTPUT)
    _rf_pi.write(ctrl_gpio, RF_SWITCH_HIGH)
    _rf_pi.set_mode(en_gpio, pigpio.OUTPUT)
    _rf_pi.write(en_gpio, RF_SWITCH_LOW)

def configure_rf_switches(pi):
    global _rf_pi
    _rf_pi = pi
    _config_HMC349(RF_SWITCH_EN_GPIO, RF_SWITCH_CTRL_GPIO)

def shutdown_rf_switches():
    _switch_HMC349(RF_SWITCH_CTRL_GPIO, RF_SWITCH_HIGH)
    _switch_HMC349(RF_SWITCH_EN_GPIO, RF_SWITCH_HIGH)

def switch_rf_switch_On():
    #print("SWITCHING ON RF SWITCH", flush=True)
    _switch_HMC349(RF_SWITCH_CTRL_GPIO, RF_SWITCH_LOW)

def switch_rf_switch_Off():
    #print("SWITCHING OFF RF SWITCH", flush=True)
    _switch_HMC349(RF_SWITCH_CTRL_GPIO, RF_SWITCH_HIGH)

*/
