package main

import (
	"fmt"
	"os"
	"strconv"
)

// Quit this will close the current tab or view that is open
func (v *View) Quit() bool {
	v.CloseBuffer()
	screen.Fini()
	os.Exit(0)

	return false
}

// UpN moves the cursor up by amount
func (v *View) UpN(amount int) bool {
	proposedY := v.Line - amount
	if proposedY < 0 {
		proposedY = 0
	} else if proposedY >= v.Buf.NumLines-1 {
		proposedY = v.Buf.NumLines - 1
	}

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

// End moves the cursor to the end of the buffer
func (v *View) Start() bool {
	v.Line = v.Buf.Start()
	return true
}

// End moves the cursor to the end of the buffer
func (v *View) End() bool {
	v.Line = v.Buf.End()
	return true
}

// Find opens a prompt and searches forward for the input
func (v *View) Find() bool {
	searchStr := ""
	BeginSearch(searchStr)
	return true
}

// FindNext searches forwards for the last used search term
func (v *View) FindNext() bool {
	searchStart = v.Line + 1
	if lastSearch == "" {
		return true
	}
	messenger.Message("Finding: " + lastSearch)
	Search(lastSearch, v, true)
	return true
}

// FindPrevious searches backwards for the last used search term
func (v *View) FindPrevious() bool {
	searchStart = v.Line
	messenger.Message("Finding: " + lastSearch)
	Search(lastSearch, v, false)
	return true
}

// ClearStatus clears the messenger bar
func (v *View) ClearStatus() bool {
	messenger.Message("")
	return false
}

func (v *View) JumpLine() bool {
	message := fmt.Sprintf("Jump to line (1 - %v) # ", v.Buf.NumLines)
	linestring, canceled := messenger.Prompt(message, "", "LineNumber", NoCompletion)
	if canceled {
		return false
	}
	lineint, err := strconv.Atoi(linestring)
	lineint = lineint - 1 // fix offset
	if err != nil {
		messenger.Error(err) // return errors
		return false
	}
	// Move cursor and view if possible.
	if lineint < v.Buf.NumLines && lineint >= 0 {
		v.Line = lineint
		return true
	}
	messenger.Error("Only ", v.Buf.NumLines, " lines to jump")
	return false
}
