/*
 *  Q-100 Receiver & Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package spClient

import (
	"context"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

const (
	config_Origin = "https://eshail.batc.org.uk/"
	config_Url    = "wss://eshail.batc.org.uk/wb/fft/fft_ea7kirsatcontroller:443/wss"
	kNumPoints    = 918
)

type (
	SpData_t struct {
		Yp          []float32
		BeaconLevel float32
	}
)

var (
	Xp = make([]float32, kNumPoints) // x coordinates from 0.0 to 100.0
)

func ReadSpectrumServer(ctx context.Context, spDataChan chan<- SpData_t) {
	var (
		ws     *websocket.Conn
		err    error
		n      int
		bytes  = make([]byte, 2048) // larger than 1844
		spData = SpData_t{
			Yp:          make([]float32, kNumPoints),
			BeaconLevel: 0.5,
		}
	)

	Xp[0] = 0
	for i := 1; i < kNumPoints-1; i++ {
		Xp[i] = 100.0 * (float32(i) / float32(kNumPoints))
	}
	Xp[kNumPoints-1] = 100

	// TODO: needs a timeout. see https://pkg.go.dev/nhooyr.io/websocket
	//	which uses: ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	const MAXTRIES = 10
	for i := 1; i <= MAXTRIES; i++ {
		log.Printf("INFO Dial attempt %v", i)
		ws, err = websocket.Dial(config_Url, "", config_Origin)
		if err == nil {
			break
		}
		if i == MAXTRIES {
			log.Fatalf("FATAL Dial Aborted after %v attemps\n", i)
		}
		time.Sleep(time.Millisecond * 500)
	}

	for {
		select {
		case <-ctx.Done():
			ws.Close()
			log.Printf("CANCEL ----- spClient has cancelled")
			return
		default:
		}

		if n, err = ws.Read(bytes); err != nil {
			// TODO: this is a PROBLEM - need to find a better way to recover
			log.Fatalf("FATAL Read failed: %v", err)
		}
		if n != 1844 {
			log.Printf("WARN reading : bytes != 1844\n")
			continue
		}

		// process the bytes
		// var count = 0
		for i := 0; i < 1836; {
			word := uint16(bytes[i]) + uint16(bytes[i+1])<<8
			// count++
			// log.Printf("INFO count = %v\n", count)
			if word < 8192 {
				word = 8192
			}
			spData.Yp[i/2] = float32(word-uint16(8192)) / float32(520) // normalize to 0 to 100
			i += 2
		}
		// log.Printf("INFO count = %v\n", count)
		spData.Yp[0] = 0
		spData.Yp[kNumPoints-1] = 0

		spData.BeaconLevel = 0
		for i := 32; i <= 133; i++ { // beacon center is 103
			spData.BeaconLevel += spData.Yp[i]
		}
		spData.BeaconLevel = spData.BeaconLevel / 103
		// log.Printf("INFO beacon level %v : Yp[i] %v", spData.BeaconLevel, spData.Yp[103])

		spDataChan <- spData

	}
}
