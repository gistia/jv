package main

import (
	"os"
	"strconv"
	"time"

	"github.com/zyedidia/tcell"
)

type ViewType struct {
	kind     int
	readonly bool // The file cannot be edited
	scratch  bool // The file cannot be saved
}

var (
	vtDefault = ViewType{0, false, false}
	vtHelp    = ViewType{1, true, true}
	vtLog     = ViewType{2, true, true}
	vtScratch = ViewType{3, false, true}
)

// The View struct stores information about a view into a buffer.
// It stores information about the cursor, and the viewport
// that the user sees the buffer from.
type View struct {
	// Y location
	Line int

	// The topmost line, used for vertical scrolling
	Topline int
	// The leftmost column, used for horizontal scrolling
	leftCol int

	// Specifies whether or not this view holds a help buffer
	Type ViewType

	// Actual width and height
	Width  int
	Height int

	LockWidth  bool
	LockHeight bool

	// Where this view is located
	x, y int

	// How much to offset because of line numbers
	lineNumOffset int

	// Holds the list of gutter messages
	messages map[string][]GutterMessage

	// This is the index of this view in the views array
	Num int
	// What tab is this view stored in
	TabNum int

	// The buffer
	Buf *Buffer
	// The statusline
	// sline Statusline

	// Since tcell doesn't differentiate between a mouse release event
	// and a mouse move event with no keys pressed, we need to keep
	// track of whether or not the mouse was pressed (or not released) last event to determine
	// mouse release events
	mouseReleased bool

	// This stores when the last click was
	// This is useful for detecting double and triple clicks
	lastClickTime time.Time
	lastLoc       Loc

	// lastCutTime stores when the last ctrl+k was issued.
	// It is used for clearing the clipboard to replace it with fresh cut lines.
	lastCutTime time.Time

	// freshClip returns true if the clipboard has never been pasted.
	freshClip bool

	// Was the last mouse event actually a double click?
	// Useful for detecting triple clicks -- if a double click is detected
	// but the last mouse event was actually a double click, it's a triple click
	doubleClick bool
	// Same here, just to keep track for mouse move events
	tripleClick bool

	cellview *CellView
}

// NewView returns a new fullscreen view
func NewView(buf *Buffer) *View {
	screenW, screenH := screen.Size()
	return NewViewWidthHeight(buf, screenW, screenH)
}

// NewViewWidthHeight returns a new view with the specified width and height
// Note that w and h are raw column and row values
func NewViewWidthHeight(buf *Buffer, w, h int) *View {
	v := new(View)

	v.x, v.y = 0, 0

	v.Width = w
	v.Height = h
	v.cellview = new(CellView)

	v.ToggleTabbar()

	v.OpenBuffer(buf)

	v.messages = make(map[string][]GutterMessage)

	// v.sline = Statusline{
	// 	view: v,
	// }

	// if v.Buf.Settings["statusline"].(bool) {
	if true {
		v.Height--
	}

	return v
}

// ToggleStatusLine creates an extra row for the statusline if necessary
func (v *View) ToggleStatusLine() {
	if v.Buf.Settings["statusline"].(bool) {
		v.Height--
	} else {
		v.Height++
	}
}

// ToggleTabbar creates an extra row for the tabbar if necessary
func (v *View) ToggleTabbar() {
	if v.y == 1 {
		v.y = 0
		v.Height++
	}
}

// ScrollUp scrolls the view up n lines (if possible)
func (v *View) ScrollUp(n int) {
	// Try to scroll by n but if it would overflow, scroll by 1
	if v.Topline-n >= 0 {
		v.Topline -= n
	} else if v.Topline > 0 {
		v.Topline--
	}
}

// ScrollDown scrolls the view down n lines (if possible)
func (v *View) ScrollDown(n int) {
	Log.Println("ScrollDown", n)
	// Try to scroll by n but if it would overflow, scroll by 1
	Log.Println("ScrollDown", v.Topline+n <= v.Buf.NumLines)
	if v.Topline+n <= v.Buf.NumLines {
		Log.Println("TopLine", n)
		v.Topline += n
	} else if v.Topline < v.Buf.NumLines-1 {
		Log.Println("TopLine++")
		v.Topline++
	}
}

// CanClose returns whether or not the view can be closed
// If there are unsaved changes, the user will be asked if the view can be closed
// causing them to lose the unsaved changes
func (v *View) CanClose() bool {
	return true
}

// OpenBuffer opens a new buffer in this view.
// This resets the topline, event handler and cursor.
func (v *View) OpenBuffer(buf *Buffer) {
	screen.Clear()
	v.CloseBuffer()
	v.Buf = buf
	v.Line = buf.Y
	v.Topline = 0
	v.leftCol = 0
	Log.Println("Relocate from OpenBuffer")
	v.Relocate()
	v.messages = make(map[string][]GutterMessage)

	// Set mouseReleased to true because we assume the mouse is not being pressed when
	// the editor is opened
	v.mouseReleased = true
	v.lastClickTime = time.Time{}
}

// Open opens the given file in the view
func (v *View) Open(filename string) {
	filename = ReplaceHome(filename)
	file, err := os.Open(filename)
	fileInfo, _ := os.Stat(filename)

	if err == nil && fileInfo.IsDir() {
		messenger.Error(filename, " is a directory")
		return
	}

	defer file.Close()

	var buf *Buffer
	if err != nil {
		messenger.Message(err.Error())
		// File does not exist -- create an empty buffer with that name
		buf = NewBufferFromString("", filename)
	} else {
		buf = NewBuffer(file, FSize(file), filename)
	}
	v.OpenBuffer(buf)
}

// CloseBuffer performs any closing functions on the buffer
func (v *View) CloseBuffer() {
}

// ReOpen reloads the current buffer
func (v *View) ReOpen() {
	if v.CanClose() {
		screen.Clear()
		// v.Buf.ReOpen()
		Log.Println("Relocate from ReOpen")
		v.Relocate()
	}
}

// HSplit opens a horizontal split with the given buffer
func (v *View) HSplit(buf *Buffer) {
	// i := 0
	// if v.Buf.Settings["splitBottom"].(bool) {
	// 	i = 1
	// }
	// v.splitNode.HSplit(buf, v.Num+i)
}

// VSplit opens a vertical split with the given buffer
func (v *View) VSplit(buf *Buffer) {
	// i := 0
	// if v.Buf.Settings["splitRight"].(bool) {
	// 	i = 1
	// }
	// v.splitNode.VSplit(buf, v.Num+i)
}

// HSplitIndex opens a horizontal split with the given buffer at the given index
func (v *View) HSplitIndex(buf *Buffer, splitIndex int) {
	// v.splitNode.HSplit(buf, splitIndex)
}

// VSplitIndex opens a vertical split with the given buffer at the given index
func (v *View) VSplitIndex(buf *Buffer, splitIndex int) {
	// v.splitNode.VSplit(buf, splitIndex)
}

func (v *View) Bottomline() int {
	return v.Topline + v.Height
}

// Relocate moves the view window so that the cursor is in view
// This is useful if the user has scrolled far away, and then starts typing
func (v *View) Relocate() bool {
	height := v.Bottomline() - v.Topline
	Log.Println("Relocate - height", height)
	ret := false
	cy := v.Line
	scrollmargin := 3
	Log.Println("cy", cy)
	Log.Println("v.Topline", v.Topline)
	Log.Println("v.Topline+scrollmargin", v.Topline+scrollmargin)
	if cy < v.Topline+scrollmargin && cy > scrollmargin-1 {
		Log.Println("realc topline")
		v.Topline = cy - scrollmargin
		ret = true
	} else if cy < v.Topline {
		Log.Println("topline y")
		v.Topline = cy
		ret = true
	}
	if cy > v.Topline+height-1-scrollmargin && cy < v.Buf.NumLines-scrollmargin {
		Log.Println("recalc topline 2")
		v.Topline = cy - height + 1 + scrollmargin
		Log.Println("recalc new Topline", v.Topline)
		ret = true
	} else if cy >= v.Buf.NumLines-scrollmargin && cy > height {
		Log.Println("set topline 2")
		v.Topline = v.Buf.NumLines - height
		ret = true
	}

	return ret
}

func (v *View) SetLine(y int) bool {
	v.Line = y
	return true
}

// HandleEvent handles an event passed by the main loop
func (v *View) HandleEvent(event tcell.Event) {
	switch e := event.(type) {
	case *tcell.EventKey:
		Log.Println("key", e.Key())
		if e.Key() == 256 || e.Key() == tcell.KeyCtrlC {
			Quit([]string{""})
		}
		if e.Key() == tcell.KeyDown {
			v.Down()
		}
		if e.Key() == tcell.KeyUp {
			v.Up()
		}
		if e.Key() == tcell.KeyPgDn {
			Log.Println("PageDown")
			v.PageDown()
		}
	}

	Log.Println("Relocate from HandleEvent")
	v.Relocate()
	v.Relocate()
}

// GutterMessage creates a message in this view's gutter
func (v *View) GutterMessage(section string, lineN int, msg string, kind int) {
	lineN--
	gutterMsg := GutterMessage{
		lineNum: lineN,
		msg:     msg,
		kind:    kind,
	}
	for _, v := range v.messages {
		for _, gmsg := range v {
			if gmsg.lineNum == lineN {
				return
			}
		}
	}
	messages := v.messages[section]
	v.messages[section] = append(messages, gutterMsg)
}

// ClearGutterMessages clears all gutter messages from a given section
func (v *View) ClearGutterMessages(section string) {
	v.messages[section] = []GutterMessage{}
}

// ClearAllGutterMessages clears all the gutter messages
func (v *View) ClearAllGutterMessages() {
	for k := range v.messages {
		v.messages[k] = []GutterMessage{}
	}
}

// Opens the given help page in a new horizontal split
func (v *View) openHelp(helpPage string) {
}

func (v *View) DisplayView() {
	// if v.Buf.Settings["softwrap"].(bool) && v.leftCol != 0 {
	// 	v.leftCol = 0
	// }

	if v.Type == vtLog {
		// Log views should always follow the cursor...
		Log.Println("Relocate from vtLog")
		v.Relocate()
	}

	// We need to know the string length of the largest line number
	// so we can pad appropriately when displaying line numbers
	maxLineNumLength := len(strconv.Itoa(v.Buf.NumLines))

	v.lineNumOffset = 0

	height := v.Height
	width := v.Width
	left := v.leftCol
	top := v.Topline

	v.cellview.Draw(v.Buf, top, height, left, width-v.lineNumOffset)

	displayLineNumber := true
	lineNumberPadding := 1
	realLineN := top - 1
	visualLineN := 0
	var line []*Char
	for visualLineN, line = range v.cellview.lines {
		screenX := 0
		realLineN++
		if displayLineNumber {
			lineNumStyle := defStyle
			lineNum := strconv.Itoa(realLineN + 1)

			// padding before
			for i := 0; i < lineNumberPadding; i++ {
				screen.SetContent(screenX, visualLineN, ' ', nil, lineNumStyle)
				screenX++
			}
			for i := 0; i < maxLineNumLength-len(lineNum); i++ {
				screen.SetContent(screenX, visualLineN, ' ', nil, lineNumStyle)
				screenX++
			}

			for _, ch := range lineNum {
				screen.SetContent(screenX, visualLineN, ch, nil, lineNumStyle)
				screenX++
			}

			// padding after
			for i := 0; i < lineNumberPadding; i++ {
				screen.SetContent(screenX, visualLineN, ' ', nil, lineNumStyle)
				screenX++
			}
		}
		lineStyle := defStyle
		if v.Line == visualLineN {
			lineStyle = defStyle.Reverse(true)
		}
		for _, ch := range line {
			// charStyle := lineStyle
			// if ch.style != nil {
			// }
			charStyle := ch.style
			if v.Line == realLineN {
				charStyle = defStyle.Reverse(true)
			}
			screen.SetContent(screenX, visualLineN, ch.drawChar, nil, charStyle)
			screenX++
		}
		for screenX < width {
			screen.SetContent(screenX, visualLineN, ' ', nil, lineStyle)
			screenX++
		}
	}
}

// Display renders the view, the cursor, and statusline
func (v *View) Display() {
	// if globalSettings["termtitle"].(bool) {
	screen.SetTitle("micro: " + v.Buf.GetName())
	// }
	v.DisplayView()
	// _, screenH := screen.Size()
	// if v.Buf.Settings["statusline"].(bool) {
	// v.sline.Display()
	// } else if (v.y + v.Height) != screenH-1 {
	for x := 0; x < v.Width; x++ {
		screen.SetContent(v.x+x, v.y+v.Height, '-', nil, defStyle.Reverse(true))
	}
	// }
}
