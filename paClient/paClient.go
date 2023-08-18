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
	"q100transmitter/mylogger"
	"strings"
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
	mylogger.Warn.Printf("SvrClient will stop... - NOT IMPLELENTED")
	// is it coonected?  send an EOF
	mylogger.Info.Printf("SvrClient has stopped - NOT IMPLELENTED")
}

// http://www.inanzzz.com/index.php/post/j3n1/creating-a-concurrent-tcp-client-and-server-example-with-golang
func readServer(cfg *SvrConfig, ch chan SvrData) {
	url := fmt.Sprintf("%s:%d", cfg.Url, cfg.Port)
	mylogger.Info.Printf("Client %v connected", url)
	conn, err := net.Dial("tcp", url)
	if err != nil {
		mylogger.Error.Printf("Failed to connect to: %v", url)
		sd := SvrData{}
		sd.Status = "Not connected"
		ch <- sd
		return
	}
	defer conn.Close()

	// clientReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

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
				if _, err = conn.Write([]byte(clientRequest + "\n")); err != nil {
					mylogger.Error.Printf("failed to send the client request: %v\n", err)
				}
			case io.EOF:
				mylogger.Info.Printf("client closed the connection")
				return
			default:
				mylogger.Error.Printf("client error: %v\n", err)
				return
			}

			// Waiting for the server response
			serverResponse, err := serverReader.ReadString('\n')

			switch err {
			case nil:
				// sd.Status = serverResponse
				sd.Status = strings.TrimSpace(serverResponse)
				ch <- sd
			case io.EOF:
				mylogger.Warn.Printf("server closed the connection")
				return
			default:
				mylogger.Warn.Printf("server error: %v\n", err)
				return
			}
		}
	}
}
