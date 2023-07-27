/*
 *  Q-100 Logger
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package logger

import (
	"log"
	"os"
)

var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Fatal *log.Logger

	logFile *os.File
)

func Open(output string) {
	logFile, err := os.OpenFile(output, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}

	flags := log.Ldate | log.Ltime | log.Lshortfile
	Info = log.New(logFile, "INFO: ", flags)
	Warn = log.New(logFile, "WARN: ", flags)
	Error = log.New(logFile, "ERROR: ", flags)
	Fatal = log.New(logFile, "FATAL: ", flags)
}

func Close() {
	logFile.Close()
}

// func Write() {

// }
