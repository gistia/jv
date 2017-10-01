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
	// Try to scroll by n but if it would overflow, scroll by 1
	if v.Topline+n <= v.Buf.NumLines {
		v.Topline += n
	} else if v.Topline < v.Buf.NumLines-1 {
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
	// if !v.Buf.Settings["softwrap"].(bool) {
	// 	return v.Topline + v.Height
	// }

	screenX, screenY := 0, 0
	numLines := 0
	for lineN := v.Topline; lineN < v.Topline+v.Height; lineN++ {
		line := v.Buf.Line(lineN)

		colN := 0
		for _, ch := range line {
			if screenX >= v.Width-v.lineNumOffset {
				screenX = 0
				screenY++
			}

			if ch == '\t' {
				screenX += int(v.Buf.Settings["tabsize"].(float64)) - 1
			}

			screenX++
			colN++
		}
		screenX = 0
		screenY++
		numLines++

		if screenY >= v.Height {
			break
		}
	}
	return numLines + v.Topline
}

// Relocate moves the view window so that the cursor is in view
// This is useful if the user has scrolled far away, and then starts typing
func (v *View) Relocate() bool {
	height := v.Bottomline() - v.Topline
	ret := false
	cy := v.Line
	// scrollmargin := int(v.Buf.Settings["scrollmargin"].(float64))
	scrollmargin := 0
	if cy < v.Topline+scrollmargin && cy > scrollmargin-1 {
		v.Topline = cy - scrollmargin
		ret = true
	} else if cy < v.Topline {
		v.Topline = cy
		ret = true
	}
	if cy > v.Topline+height-1-scrollmargin && cy < v.Buf.NumLines-scrollmargin {
		v.Topline = cy - height + 1 + scrollmargin
		ret = true
	} else if cy >= v.Buf.NumLines-scrollmargin && cy > height {
		v.Topline = v.Buf.NumLines - height
		ret = true
	}

	// if !v.Buf.Settings["softwrap"].(bool) {
	// 	cx := v.Cursor.GetVisualX()
	// 	if cx < v.leftCol {
	// 		v.leftCol = cx
	// 		ret = true
	// 	}
	// 	if cx+v.lineNumOffset+1 > v.leftCol+v.Width {
	// 		v.leftCol = cx - v.Width + v.lineNumOffset + 1
	// 		ret = true
	// 	}
	// }
	return ret
}

// MoveToMouseClick moves the cursor to location x, y assuming x, y were given
// by a mouse click
func (v *View) MoveToMouseClick(x, y int) {
	// if y-v.Topline > v.Height-1 {
	// 	v.ScrollDown(1)
	// 	y = v.Height + v.Topline - 1
	// }
	// if y < 0 {
	// 	y = 0
	// }
	// if x < 0 {
	// 	x = 0
	// }

	// x, y = v.GetSoftWrapLocation(x, y)
	// // x = v.Cursor.GetCharPosInLine(y, x)
	// if x > Count(v.Buf.Line(y)) {
	// 	x = Count(v.Buf.Line(y))
	// }
	// v.Cursor.X = x
	// v.Cursor.Y = y
	// v.Cursor.LastVisualX = v.Cursor.GetVisualX()
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
	}

	// // This bool determines whether the view is relocated at the end of the function
	// // By default it's true because most events should cause a relocate
	// relocate := true

	// switch e := event.(type) {
	// case *tcell.EventKey:
	// 	// Check first if input is a key binding, if it is we 'eat' the input and don't insert a rune
	// 	isBinding := false
	// 	for key, actions := range bindings {
	// 		if e.Key() == key.keyCode {
	// 			if e.Key() == tcell.KeyRune {
	// 				if e.Rune() != key.r {
	// 					continue
	// 				}
	// 			}
	// 			if e.Modifiers() == key.modifiers {
	// 				for _, c := range v.Buf.cursors {
	// 					ok := v.SetCursor(c)
	// 					if !ok {
	// 						break
	// 					}
	// 					relocate = false
	// 					isBinding = true
	// 					relocate = v.ExecuteActions(actions) || relocate
	// 				}
	// 				v.SetCursor(&v.Buf.Cursor)
	// 				v.Buf.MergeCursors()
	// 				break
	// 			}
	// 		}
	// 	}
	// 	if !isBinding && e.Key() == tcell.KeyRune {
	// 		// Check viewtype if readonly don't insert a rune (readonly help and log view etc.)
	// 		if v.Type.readonly == false {
	// 			for _, c := range v.Buf.cursors {
	// 				v.SetCursor(c)

	// 				// Insert a character
	// 				if v.Cursor.HasSelection() {
	// 					v.Cursor.DeleteSelection()
	// 					v.Cursor.ResetSelection()
	// 				}
	// 				v.Buf.Insert(v.Cursor.Loc, string(e.Rune()))

	// 				for pl := range loadedPlugins {
	// 					_, err := Call(pl+".onRune", string(e.Rune()), v)
	// 					if err != nil && !strings.HasPrefix(err.Error(), "function does not exist") {
	// 						TermMessage(err)
	// 					}
	// 				}

	// 				if recordingMacro {
	// 					curMacro = append(curMacro, e.Rune())
	// 				}
	// 			}
	// 			v.SetCursor(&v.Buf.Cursor)
	// 		}
	// 	}
	// case *tcell.EventPaste:
	// 	// Check viewtype if readonly don't paste (readonly help and log view etc.)
	// 	if v.Type.readonly == false {
	// 		if !PreActionCall("Paste", v) {
	// 			break
	// 		}

	// 		for _, c := range v.Buf.cursors {
	// 			v.SetCursor(c)
	// 			v.paste(e.Text())
	// 		}
	// 		v.SetCursor(&v.Buf.Cursor)

	// 		PostActionCall("Paste", v)
	// 	}
	// case *tcell.EventMouse:
	// 	// Don't relocate for mouse events
	// 	relocate = false

	// 	button := e.Buttons()

	// 	for key, actions := range bindings {
	// 		if button == key.buttons && e.Modifiers() == key.modifiers {
	// 			for _, c := range v.Buf.cursors {
	// 				ok := v.SetCursor(c)
	// 				if !ok {
	// 					break
	// 				}
	// 				relocate = v.ExecuteActions(actions) || relocate
	// 			}
	// 			v.SetCursor(&v.Buf.Cursor)
	// 			v.Buf.MergeCursors()
	// 		}
	// 	}

	// 	for key, actions := range mouseBindings {
	// 		if button == key.buttons && e.Modifiers() == key.modifiers {
	// 			for _, action := range actions {
	// 				action(v, true, e)
	// 			}
	// 		}
	// 	}

	// 	switch button {
	// 	case tcell.ButtonNone:
	// 		// Mouse event with no click
	// 		if !v.mouseReleased {
	// 			// Mouse was just released

	// 			x, y := e.Position()
	// 			x -= v.lineNumOffset - v.leftCol + v.x
	// 			y += v.Topline - v.y

	// 			// Relocating here isn't really necessary because the cursor will
	// 			// be in the right place from the last mouse event
	// 			// However, if we are running in a terminal that doesn't support mouse motion
	// 			// events, this still allows the user to make selections, except only after they
	// 			// release the mouse

	// 			if !v.doubleClick && !v.tripleClick {
	// 				v.MoveToMouseClick(x, y)
	// 				v.Cursor.SetSelectionEnd(v.Cursor.Loc)
	// 				v.Cursor.CopySelection("primary")
	// 			}
	// 			v.mouseReleased = true
	// 		}
	// 	}
	// }

	// if relocate {
	// 	v.Relocate()
	// 	// We run relocate again because there's a bug with relocating with softwrap
	// 	// when for example you jump to the bottom of the buffer and it tries to
	// 	// calculate where to put the topline so that the bottom line is at the bottom
	// 	// of the terminal and it runs into problems with visual lines vs real lines.
	// 	// This is (hopefully) a temporary solution
	// 	v.Relocate()
	// }
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
		v.Relocate()
	}

	// We need to know the string length of the largest line number
	// so we can pad appropriately when displaying line numbers
	maxLineNumLength := len(strconv.Itoa(v.Buf.NumLines))

	// if v.Buf.Settings["ruler"] == true {
	// 	// + 1 for the little space after the line number
	// 	v.lineNumOffset = maxLineNumLength + 1
	// } else {
	v.lineNumOffset = 0
	// }

	// We need to add to the line offset if there are gutter messages
	var hasGutterMessages bool
	for _, v := range v.messages {
		if len(v) > 0 {
			hasGutterMessages = true
		}
	}
	if hasGutterMessages {
		v.lineNumOffset += 2
	}

	// divider := 0
	// if v.x != 0 {
	// 	// One space for the extra split divider
	// 	v.lineNumOffset++
	// 	divider = 1
	// }

	// xOffset := v.x + v.lineNumOffset
	// yOffset := v.y

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
			screen.SetContent(screenX, visualLineN, ch.drawChar, nil, lineStyle)
			screenX++
		}
		for screenX < width {
			screen.SetContent(screenX, visualLineN, ' ', nil, lineStyle)
			screenX++
		}
		realLineN++
	}
}

// Display renders the view, the cursor, and statusline
func (v *View) Display() {
	// if globalSettings["termtitle"].(bool) {
	// 	screen.SetTitle("micro: " + v.Buf.GetName())
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
