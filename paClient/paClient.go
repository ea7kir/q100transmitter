/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package paClient

import (
	"bufio"
	"io"
	"log"
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
	// go readServer(cfg, ch)
	go new_readServer(cfg, ch)
}

// API
func Stop() {
	logger.Warn.Printf("SvrClient will stop... - NOT IMPLELENTED")
	//
	logger.Info.Printf("SvrClient has stopped - NOT IMPLELENTED")
}

// func readServer(cfg *SvrConfig, ch chan SvrData) {

// 	sd := SvrData{}
// 	count := 0

// 	for {
// 		time.Sleep(time.Second)
// 		count++
// 		str := fmt.Sprintf("SUCCESS # %v", count)
// 		sd.Status = str
// 		ch <- sd
// 	}

// }

// http://www.inanzzz.com/index.php/post/j3n1/creating-a-concurrent-tcp-client-and-server-example-with-golang
func new_readServer(cfg *SvrConfig, ch chan SvrData) {
	con, err := net.Dial("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()

	// clientReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(con)

	sd := SvrData{}
	for {
		time.Sleep(time.Second)

		for {
			// Waiting for the client request
			// clientRequest, err := clientReader.ReadString('\n')
			time.Sleep(2 * time.Second)

			switch err {
			case nil:
				clientRequest := ""
				if _, err = con.Write([]byte(clientRequest + "\n")); err != nil {
					log.Printf("failed to send the client request: %v\n", err)
				}
			case io.EOF:
				log.Println("client closed the connection")
				return
			default:
				log.Printf("client error: %v\n", err)
				return
			}

			// Waiting for the server response
			serverResponse, err := serverReader.ReadString('\n')

			switch err {
			case nil:
				sd.Status = serverResponse
				ch <- sd
			case io.EOF:
				log.Println("server closed the connection")
				return
			default:
				log.Printf("server error: %v\n", err)
				return
			}
		}
	}
}
