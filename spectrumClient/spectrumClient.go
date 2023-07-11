/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package spectrumClient

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"q100transmitter/logger"
	"time"

	"github.com/gorilla/websocket"
)

// API
type (
	SpConfig struct {
		Url  string
		Port int16
	}
	SpData struct {
		Yp                        []float32
		BeaconLevel               float32
		MarkerCentre, MarkerWidth float32
	}
)

// API
var (
	Xp = make([]float32, numPoints) // x coordinates from 0.0 to 100.0
)

// API
func Intitialize(cfg *SpConfig, ch chan SpData) {
	spChannel = ch
	Xp[0] = 0
	for i := 1; i < numPoints-1; i++ {
		Xp[i] = 100.0 * (float32(i) / float32(numPoints))
	}
	Xp[numPoints-1] = 100
	url := fmt.Sprintf("%v:%v/", cfg.Url, cfg.Port)
	go readAndDecode(url, spChannel)
}

// API
func Stop() {
	logger.Warn.Printf("Spectrum will stop... - NOT IMPLELENTED")
	//
	logger.Info.Printf("Spectrum has stopped - NOT IMPLELENTED")
}

// Sets the spData Marker values
func SetMarker(frequency, symbolRate string) {
	spData.MarkerCentre, spData.MarkerWidth = getMarkers(frequency, symbolRate)
}

// room for 916 datapoints + start and end zero points to close the polygon
const numPoints = 918

var (
	spData = SpData{
		Yp:           make([]float32, numPoints),
		BeaconLevel:  0.5,
		MarkerCentre: 0.5,
		MarkerWidth:  0.5,
	}
	spChannel chan SpData
)

// func readCalibrationData(ch chan SpData) {
// 	logger.Info.Printf("Spectrun calibration running...")
// 	for {
// 		spData.Yp[0] = 0
// 		for i := 1; i < numPoints-2; i++ {
// 			spData.Yp[i] = rand.Float32() * 50.0
// 		}
// 		spData.Yp[numPoints-1] = 0
// 		spData.BeaconLevel = rand.Float32() * 100
// 		ch <- spData
// 		time.Sleep(3 * time.Millisecond)
// 	}
// }

func readAndDecode(url string, ch chan SpData) {
	ctx := context.Background()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	dialer := websocket.Dialer{
		//Subprotocols: []string{"json"},
	}
	c, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		logger.Fatal.Printf("Dial failed: %#v with %v\n", err, url)
		// log.Panic()
	}
	defer c.Close()
	// logger.Info.Printf("negotiated protocol: %q\n", c.Subprotocol())

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		// count := 0
		for {
			_, bytes, err := c.ReadMessage()
			if err != nil {
				// TODO: this can happed on Control-C
				logger.Warn.Printf("reading batc bytes: %#v\n", err)
				return
			}
			// message is of Type: []unit8
			if len(bytes) != 1844 {
				logger.Warn.Printf("reading : bytes != 1844\n")
				continue
			}
			// logger.Info.Printf("received length: %v, Type: %T, %d == %d", len(message), message, message[0], i)

			// begin processing the bytes
			// count = 0
			for i := 0; i < 1836; {
				word := uint16(bytes[i]) + uint16(bytes[i+1])<<8
				// count++
				// logger.Info.Printf("count = %v\n", count)
				if word < 8192 {
					word = 8192
				}
				// spData.Yp[i/2] = float32(word-uint16(8192)) / float32(52000)
				spData.Yp[i/2] = float32(word-uint16(8192)) / float32(520) // normalize to 0 to 100
				// spData.Yp[i/2] = 50.0
				i += 2
			}
			// logger.Info.Printf("count = %v\n", count)
			spData.Yp[0] = 0
			spData.Yp[numPoints-1] = 0

			spData.BeaconLevel = 0
			for i := 32; i <= 133; i++ { // beacon center is 103
				spData.BeaconLevel += spData.Yp[i]
			}
			spData.BeaconLevel = spData.BeaconLevel / 103
			// logger.Info.Printf("beacon level %v : Yp[i] %v", spData.BeaconLevel, spData.Yp[103])

			ch <- spData
		}
	}()

	//I don't know why the following is needed - yet

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			if err := c.WriteMessage(websocket.TextMessage, []byte(t.String())); err != nil {
				logger.Warn.Printf("writing: %#v\n", err)
				return
			}
		case <-interrupt:
			logger.Info.Printf("interrupting")
			if err := c.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(
					websocket.CloseNormalClosure, "",
				)); err != nil {
				logger.Warn.Printf("error closing: %#v", err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
			return
		}
	}
}
