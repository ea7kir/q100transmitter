package server

import (
	"fmt"
	"time"
)

func readServer(cfg SvrConfig, ch chan SvrData) {

	sd := SvrData{}
	count := 0

	for {
		time.Sleep(time.Second)
		count++
		str := fmt.Sprintf("SUCCESS # %v", count)
		sd.Status = str
		ch <- sd
	}

}
