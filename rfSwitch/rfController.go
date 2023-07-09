package rfSwitch

import "q100transmitter/logger"

var (
	state bool
)

func setPtt(on bool) {
	if on {
		state = true
	} else {
		state = false
	}
	logger.Info.Printf("PTT is %v", state)
}
