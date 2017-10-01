package main

import (
	"log"
	"os"
)

var (
	Log *log.Logger
)

func NewLog(logfile string) {
	file, err := os.Create(logfile)
	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
}
