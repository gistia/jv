package main

import (
	"log"

	"github.com/gistia/jv/ui"
)

func doUI(log string, logger *log.Logger) error {
	app := ui.NewApp(log)
	app.SetLogger(logger)

	app.Run()
	return nil
}
