/*
 *  Q-100 Transmitter
 *  Copyright (c) 2023 Michael Naylor EA7KIR (https://michaelnaylor.es)
 */

package logger

import (
	"log"
	"os"
)

// API
var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Fatal *log.Logger
)

// type aggregatedLogger struct {
// 	infoLogger  *log.Logger
// 	warnLogger  *log.Logger
// 	errorLogger *log.Logger
// }

func init() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	// flags := log.LstdFlags | log.Lshortfile
	flags := log.Ltime | log.Llongfile
	Info = log.New(os.Stderr, "INFO: ", flags)
	Warn = log.New(os.Stderr, "WARN: ", flags)
	Error = log.New(os.Stderr, "ERROR: ", flags)
	Fatal = log.New(os.Stderr, "FATAL: ", flags)

	// 	al := aggregatedLogger{
	// 		infoLogger:  log.New(os.Stderr, prefix:"INFO:", log.LstdFlags|log.Lshortfile),
	// 		warnLogger:  log.New(os.Stderr, "WARNING:", log.LstdFlags|log.Lshortfile),
	// 		errorLogger: log.New(os.Stderr, "ERROR:", log.LstdFlags|log.Lshortfile),
	// 	}
}

// func (l *aggregatedLogger) Info(v ...interface{}) {
// 	l.infoLogger.Printf(v...)
// }

// func (l *aggregatedLogger) Warn(v ...interface{}) {
// 	l.warnLogger.Printf(v...)
// }

//	func (l *aggregatedLogger) Error(v ...interface{}) {
//		l.errorLogger.Printf(v...)
// }
