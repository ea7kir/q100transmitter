/*
 *  Q-100 PA Server
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package logger

import (
	"log"
	"os"
)

var (
	infoLevel  *log.Logger
	warnLevel  *log.Logger
	errorLevel *log.Logger
	fatalLevel *log.Logger
)

func init() {
	flags := log.Ltime | log.Llongfile
	infoLevel = log.New(os.Stderr, "INFO: ", flags)
	warnLevel = log.New(os.Stderr, "WARN: ", flags)
	errorLevel = log.New(os.Stderr, "ERROR: ", flags)
	fatalLevel = log.New(os.Stderr, "FATAL: ", flags)
}

func Info(format string, v ...interface{}) {
	infoLevel.Printf(format, v...)
}

func Warn(format string, v ...interface{}) {
	warnLevel.Printf(format, v...)
}

func Error(format string, v ...interface{}) {
	errorLevel.Printf(format, v...)
}

func Fatal(format string, v ...interface{}) {
	fatalLevel.Fatalf(format, v...)
}
