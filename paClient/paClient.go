/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package paClient

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"q100transmitter/logger"
	"time"
)

type (
	// API
	SvrConfig struct {
		Url  string
		Port int
	}
	// API
	SvrData struct {
		Status string
	}
)

// API
func Initialize(cfg *SvrConfig, ch chan SvrData) {
	go readServer(cfg, ch)
	// TODO: create the connection in loop to retry
	// defer
	// if ok the start a tickker and go client
}

// API
func Stop() {
	logger.Warn.Printf("SvrClient will stop... - NOT IMPLELENTED")
	// is it coonected?  send an EOF
	logger.Info.Printf("SvrClient has stopped - NOT IMPLELENTED")
}

// http://www.inanzzz.com/index.php/post/j3n1/creating-a-concurrent-tcp-client-and-server-example-with-golang
func readServer(cfg *SvrConfig, ch chan SvrData) {
	url := fmt.Sprintf("%s:%d", cfg.Url, cfg.Port)
	logger.Info.Printf(">%v<\n", url)
	con, err := net.Dial("tcp", url)
	if err != nil {
		logger.Error.Printf("Failed to connect to: %v", url)
		sd := SvrData{}
		sd.Status = "Not connected"
		ch <- sd
		return
	}
	defer con.Close()

	// clientReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(con)

	sd := SvrData{}
	for {
		// TODO: better to use a ticker
		// t := time.NewTicker(2 * time.Second)
		// <-t.C
		// send request

		time.Sleep(time.Second)

		for {
			// Waiting for the client request
			// clientRequest, err := clientReader.ReadString('\n')
			time.Sleep(2 * time.Second)

			switch err {
			case nil:
				clientRequest := ""
				if _, err = con.Write([]byte(clientRequest + "\n")); err != nil {
					logger.Error.Printf("failed to send the client request: %v\n", err)
				}
			case io.EOF:
				logger.Info.Printf("client closed the connection")
				return
			default:
				logger.Error.Printf("client error: %v\n", err)
				return
			}

			// Waiting for the server response
			serverResponse, err := serverReader.ReadString('\n')

			switch err {
			case nil:
				sd.Status = serverResponse
				ch <- sd
			case io.EOF:
				logger.Warn.Printf("server closed the connection")
				return
			default:
				logger.Warn.Printf("server error: %v\n", err)
				return
			}
		}
	}
}
