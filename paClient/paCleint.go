/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package paClient

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

const (
	config_PaUrl  = "paserver.local"
	config_PaPort = 9999 //8765,
)

type (
	SvrConfig_t struct {
		Url  string
		Port int
	}
	SvrData_t struct {
		Status string
	}
)

// TODO: need to add a timeout

// http://www.inanzzz.com/index.php/post/j3n1/creating-a-concurrent-tcp-client-and-server-example-with-golang

func ReadPaServer(ctx context.Context, ch chan SvrData_t) {
	sd := SvrData_t{}
	// sd.Status = "Connecting..."
	// ch <- sd

	url := fmt.Sprintf("%s:%d", config_PaUrl, config_PaPort)
	log.Printf("INFO connecting to %v", url)

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
			// log.Fatalf("FATAL   Dial Aborted after %v attemps\n", i)
			log.Printf("ERROR Dial Aborted after %v attemps\n", i)
			// sd := SvrData_t{}
			sd.Status = "Not connected"
			ch <- sd
			return
		}
	}

	serverReader := bufio.NewReader(conn)
	clientRequest := "\n"

	ticker := time.NewTicker(time.Second * 2)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			conn.Close()
			log.Printf(("CANCEL ----- paClient has cancelled"))
			return
		case <-ticker.C:
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
				log.Printf("WARN server closed the connection")
				sd.Status = "Server closed the connection"
				ch <- sd
				return
			default:
				log.Printf("ERROR server error: %v\n", err)
				sd.Status = "Server error occured"
				ch <- sd
				return
			}
		}
	}
}
