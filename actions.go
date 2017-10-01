package main

import "os"

// Quit this will close the current tab or view that is open
func (v *View) Quit(usePlugin bool) bool {
	v.CloseBuffer()
	screen.Fini()
	os.Exit(0)

	return false
}
