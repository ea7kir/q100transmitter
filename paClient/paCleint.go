/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package paClient

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
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

var (
	done bool
)

// API
func Initialize(cfg SvrConfig, ch chan SvrData) {
	go readServer(cfg, ch)
	// TODO: create the connection in loop to retry
	// defer
	// if ok the start a tickker and go client
}

// API
func Stop() {
	log.Printf("WARN  paClient will stop... - NOT IMPLELENTED")
	// is it coonected?  send an EOF
	done = true
}

// TODO: need to add a timeout

// http://www.inanzzz.com/index.php/post/j3n1/creating-a-concurrent-tcp-client-and-server-example-with-golang
func readServer(cfg SvrConfig, ch chan SvrData) {
	sd := SvrData{}
	sd.Status = "Connecting..."
	ch <- sd

	url := fmt.Sprintf("%s:%d", cfg.Url, cfg.Port)
	log.Printf("INFO Client %v connected", url)

	const MAXTRIES = 10
	var conn net.Conn

	for i := 1; i <= MAXTRIES; i++ {
		log.Printf("INFO Dial attempt %v", i)
		new_conn, err := net.Dial("tcp", url)
		if err == nil {
			conn = new_conn
			break
		}
		if i == MAXTRIES {
			// qLog.Fatal("Dial Aborted after %v attemps\n", i)
			log.Printf("ERROR Dial Aborted after %v attemps\n", i)
			// sd := SvrData{}
			sd.Status = "Not connected"
			ch <- sd
			return
		}

	}
	defer conn.Close()

	serverReader := bufio.NewReader(conn)
	clientRequest := "\n"

	// for {

	// TODO: better to use a ticker
	// t := time.NewTicker(2 * time.Second)
	// <-t.C
	// send request

	for {
		if done {
			log.Printf("WARN  paClient will stop...")
			return
		}

		time.Sleep(2 * time.Second)

		if _, err := conn.Write([]byte(clientRequest)); err != nil {
			log.Printf("ERROR failed to send the client request: %v\n", err)
			sd.Status = "Failed to send request"
			ch <- sd
		}

		// Waiting for the server response
		serverResponse, err := serverReader.ReadString('\n')
		switch err {
		case nil:
			sd.Status = strings.TrimSpace(serverResponse)
			ch <- sd
		case io.EOF:
			log.Printf("WARN  server closed the connection")
			return
		default:
			log.Printf("ERROR server error: %v\n", err)
			return
		}
	}
}
