package logic

import (
	"gvm/core/flags"
	"log"
)

func DebugPrintln(message string) {
	if !flags.DebugEnabled {
		return
	}
	log.Println(message)
}

func DebugPrintf(format string, a ...interface{}) {
	if !flags.DebugEnabled {
		return
	}
	log.Printf(format, a...)
}
