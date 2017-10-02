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
	Log.Println("UpN - Up by", amount)
	proposedY := v.Line - amount
	Log.Println("UpN - Proposed", proposedY)
	Log.Println("UpN - NumLines", v.Buf.NumLines)
	if proposedY < 0 {
		proposedY = 0
	} else if proposedY >= v.Buf.NumLines {
		proposedY = v.Buf.NumLines
	}
	Log.Println("UpN - Actual", proposedY)

	v.Line = proposedY
}

// Down moves the cursor one line down
func (v *View) Down() {
	v.UpN(-1)
}

// Up moves the cursor one line down
func (v *View) Up() {
	v.UpN(1)
}

// PageDown scrolls the view down a page
func (v *View) PageDown() {
	Log.Println("PageDown", v.Buf.NumLines-(v.Topline+v.Height), v.Height)
	if v.Buf.NumLines-(v.Topline+v.Height) > v.Height {
		v.ScrollDown(v.Height)
	} else if v.Buf.NumLines >= v.Height {
		Log.Println("TopLine", v.Buf.NumLines-v.Height)
		v.Topline = v.Buf.NumLines - v.Height
	}
}
