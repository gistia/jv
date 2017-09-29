package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	logfile := "/Users/fcoury/logs/jvg.log"
	var dlog *log.Logger
	if logfile != "" {
		f, e := os.Create(logfile)
		if e == nil {
			dlog = log.New(f, "DEBUG:", log.LstdFlags)
			log.SetOutput(f)
		}
	}
	if e := doUI("/Users/fcoury/logs/large.log", dlog); e != nil {
		fmt.Println("worked")
	} else {
		return
	}
}
