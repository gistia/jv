package main

import "os"

// Quit this will close the current tab or view that is open
func (v *View) Quit() bool {
	v.CloseBuffer()
	screen.Fini()
	os.Exit(0)

	return false
}

// UpN moves the cursor up by amount
func (v *View) UpN(amount int) bool {
	Log.Println("UpN - Up by", amount)
	proposedY := v.Line - amount
	Log.Println("UpN - Proposed", proposedY)
	Log.Println("UpN - NumLines", v.Buf.NumLines)
	if proposedY < 0 {
		proposedY = 0
	} else if proposedY >= v.Buf.NumLines-1 {
		proposedY = v.Buf.NumLines - 1
	}
	Log.Println("UpN - Actual", proposedY)

	v.Line = proposedY
	return true
}

func (v *View) DownN(amount int) bool {
	return v.UpN(-amount)
}

// Down moves the cursor one line down
func (v *View) Down() bool {
	v.UpN(-1)
	return true
}

// Up moves the cursor one line down
func (v *View) Up() bool {
	v.UpN(1)
	return true
}

// PageUp scrolls the view up a page
func (v *View) PageUp() bool {
	return v.UpN(v.Height)
}

// PageDown scrolls the view down a page
func (v *View) PageDown() bool {
	return v.DownN(v.Height)
}
