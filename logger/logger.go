/*
 *  Q-100 PA Transmitter
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
)

func init() {
	flags := log.Ltime | log.Lshortfile
	Info = log.New(os.Stderr, "INFO: ", flags)
	Warn = log.New(os.Stderr, "WARN: ", flags)
	Error = log.New(os.Stderr, "ERROR: ", flags)
	Fatal = log.New(os.Stderr, "FATAL: ", flags)
}
