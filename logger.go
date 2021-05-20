package closer

import "log"

// ----------------------
//   Logger
// ----------------------

type Logger interface {
	Error(v ...interface{})
	// Debug(v ...interface{})
	Info(v ...interface{})
}

type defaultLogger struct{}

var logger = defaultLogger{}

func DefaultLogger() Logger {
	return logger
}

//
func (d defaultLogger) Error(v ...interface{}) {
	log.Printf("Err: %v \n", v)
}

//
/*
func (d defaultLogger) Debug(v ...interface{}) {
	log.Printf("Debug: %v \n", v)
}
*/

//
func (d defaultLogger) Info(v ...interface{}) {
	log.Printf("Info: %v \n", v)
}
