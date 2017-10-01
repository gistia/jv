package main

import "os"

// Quit this will close the current tab or view that is open
func (v *View) Quit(usePlugin bool) bool {
	v.CloseBuffer()
	screen.Fini()
	os.Exit(0)

	return false
}

// UpN moves the cursor up by amount
func (v *View) UpN(amount int) {
	Log.Println("Up by", amount)
	proposedY := v.Line + amount
	Log.Println("Proposed", proposedY)
	if proposedY < 0 {
		proposedY = 0
	} else if proposedY >= v.Buf.NumLines {
		proposedY = v.Buf.NumLines - 1
	}

	v.Line = proposedY
}

// Down moves the cursor one line down
func (v *View) Down() {
	v.UpN(1)
}

// Up moves the cursor one line down
func (v *View) Up() {
	v.UpN(-1)
}
