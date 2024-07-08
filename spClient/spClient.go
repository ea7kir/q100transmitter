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

type (
	SpConfig struct {
		Url    string
		Origin string
	}
	SpData struct {
		Yp           []float32
		BeaconLevel  float32
		MarkerCentre float32
		MarkerWidth  float32
	}
)

var (
	Xp = make([]float32, numPoints) // x coordinates from 0.0 to 100.0
)

func Start(ctx context.Context, cfg SpConfig, ch chan SpData) {
	Xp[0] = 0
	for i := 1; i < numPoints-1; i++ {
		Xp[i] = 100.0 * (float32(i) / float32(numPoints))
	}
	Xp[numPoints-1] = 100

	go readAndDecode(ctx, cfg, ch)
}

// END API *******************************************************

// room for 916 datapoints + start and end zero points to close the polygon
const numPoints = 918

var (
	spData = SpData{
		Yp:           make([]float32, numPoints),
		BeaconLevel:  0.5,
		MarkerCentre: 0.5,
		MarkerWidth:  0.5,
	}
)

// TODO: needs a timeout. see https://pkg.go.dev/nhooyr.io/websocket
//	which uses: ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

// forever go routine called from Start
func readAndDecode(ctx context.Context, cfg SpConfig, ch chan SpData) {
	const MAXTRIES = 10
	var ws *websocket.Conn

	for i := 1; i <= MAXTRIES; i++ {
		log.Printf("INFO Dial attempt %v", i)
		new_ws, err := websocket.Dial(cfg.Url, "", cfg.Origin)
		if err == nil {
			ws = new_ws
			break
		}
		if i == MAXTRIES {
			log.Fatalf("FATAL   Dial Aborted after %v attemps\n", i)
		}
		time.Sleep(time.Millisecond * 500)
	}

	var bytes = make([]byte, 2048) // larger than 1844
	var n int
	var err error

	for {
		if ctx.Err() != nil {
			time.Sleep(time.Duration(time.Second))
			ws.Close()
			log.Printf("INFO ----- Cancelled readAndDecode and ws closed")
			return
		}
		if n, err = ws.Read(bytes); err != nil {
			log.Printf("WARN  Read failed: %v", err)
			continue
		}
		if n != 1844 {
			log.Printf("WARN  reading : bytes != 1844\n")
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
		spData.Yp[numPoints-1] = 0

		spData.BeaconLevel = 0
		for i := 32; i <= 133; i++ { // beacon center is 103
			spData.BeaconLevel += spData.Yp[i]
		}
		spData.BeaconLevel = spData.BeaconLevel / 103
		// log.Printf("INFO beacon level %v : Yp[i] %v", spData.BeaconLevel, spData.Yp[103])

		ch <- spData
	}
}
