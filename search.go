package main

import (
	"regexp"

	"github.com/zyedidia/tcell"
)

var (
	// What was the last search
	lastSearch string

	// Where should we start the search down from (or up from)
	searchStart int

	// Is there currently a search in progress
	searching bool

	// Stores the history for searching
	searchHistory []string
)

// BeginSearch starts a search
func BeginSearch(searchStr string) {
	searchHistory = append(searchHistory, "")
	messenger.historyNum = len(searchHistory) - 1
	searching = true
	messenger.response = searchStr
	messenger.cursorx = Count(searchStr)
	messenger.Message("Find: ")
	messenger.hasPrompt = true
}

// EndSearch stops the current search
func EndSearch() {
	searchHistory[len(searchHistory)-1] = messenger.response
	searching = false
	messenger.hasPrompt = false
	messenger.Clear()
	messenger.Reset()
	if lastSearch != "" {
		messenger.Message("N Previous n Next")
	}
}

// ExitSearch exits the search mode, reset active search phrase, and clear status bar
func ExitSearch(v *View) {
	lastSearch = ""
	searching = false
	messenger.hasPrompt = false
	messenger.Clear()
	messenger.Reset()
}

// HandleSearchEvent takes an event and a view and will do a real time match from the messenger's output
// to the current buffer. It searches down the buffer.
func HandleSearchEvent(event tcell.Event, v *View) {
	switch e := event.(type) {
	case *tcell.EventKey:
		switch e.Key() {
		case tcell.KeyEscape:
			// Exit the search mode
			ExitSearch(v)
			return
		case tcell.KeyCtrlQ, tcell.KeyCtrlC, tcell.KeyEnter:
			// Done
			EndSearch()
			return
		}
	}

	messenger.HandleEvent(event, searchHistory)

	if messenger.cursorx < 0 {
		// Done
		EndSearch()
		return
	}

	if messenger.response == "" {
		// v.Cursor.ResetSelection()
		// We don't end the search though
		return
	}

	Search(messenger.response, v, true)

	v.Relocate()

	return
}

func searchDown(r *regexp.Regexp, v *View, startY, endY int) bool {
	for i := startY; i <= endY; i++ {
		var l []byte
		if i == startY {
			runes := []rune(string(v.Buf.lines[i].data))
			l = []byte(string(runes[0:]))
		} else {
			l = v.Buf.lines[i].data
		}

		match := r.FindIndex(l)

		if match != nil {
			v.Line = i
			return true
		}
	}
	return false
}

func searchUp(r *regexp.Regexp, v *View, startY, endY int) bool {
	for i := startY; i >= endY; i-- {
		var l []byte
		if i == startY {
			runes := []rune(string(v.Buf.lines[i].data))
			l = []byte(string(runes[:0]))
		} else {
			l = v.Buf.lines[i].data
		}

		match := r.FindIndex(l)

		if match != nil {
			v.Line = i
			return true
		}
	}
	return false
}

// Search searches in the view for the given regex. The down bool
// specifies whether it should search down from the searchStart position
// or up from there
func Search(searchStr string, v *View, down bool) {
	if searchStr == "" {
		return
	}
	// r, err := regexp.Compile(searchStr)
	// if v.Buf.Settings["ignorecase"].(bool) {
	r, err := regexp.Compile("(?i)" + searchStr)
	// }
	if err != nil {
		return
	}

	var found bool
	if down {
		found = searchDown(r, v, searchStart, v.Buf.End())
		if !found {
			found = searchDown(r, v, v.Buf.Start(), searchStart)
		}
	} else {
		found = searchUp(r, v, searchStart, v.Buf.Start())
		if !found {
			found = searchUp(r, v, v.Buf.End(), searchStart)
		}
	}
	if found {
		lastSearch = searchStr
	}
}
