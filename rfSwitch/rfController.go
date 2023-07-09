package rfSwitch

import "q100transmitter/logger"

var (
	state bool
)

func setPtt(on bool) bool {
	if on {
		state = true
	} else {
		state = false
	}
	logger.Info.Printf("PTT is %v", state)
	return state
}
