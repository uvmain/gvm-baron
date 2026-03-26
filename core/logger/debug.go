package logger

import (
	"log"
)

var DebugEnabled bool

func DebugPrintln(message string) {
	if !DebugEnabled {
		return
	}
	log.Println(message)
}

func DebugPrintf(format string, a ...interface{}) {
	if !DebugEnabled {
		return
	}
	log.Printf(format, a...)
}
