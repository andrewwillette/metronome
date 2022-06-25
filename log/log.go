package log

import (
	"fmt"
	"log"
	"os"
)

var debug = false

func ConfigureLog(logFile string, debugVal bool) {
	debug = debugVal
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	log.SetOutput(f)
}

func Lg(output string) {
	if debug {
		log.Println(output)
	}
}
