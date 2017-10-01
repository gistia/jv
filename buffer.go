package main

import (
	"crypto/md5"
	"io"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

var (
	// 0 - no line type detected
	// 1 - lf detected
	// 2 - crlf detected
	fileformat = 0
)

// Buffer stores the text for files that are loaded into the text editor
// It uses a rope to efficiently store the string and contains some
// simple functions for saving and wrapper functions for modifying the rope
type Buffer struct {
	// This stores all the text in the buffer as an array of lines
	*LineArray

	// Path to the file on disk
	Path string
	// Absolute path to the file on disk
	AbsPath string
	// Name of the buffer on the status line
	name string

	NumLines int

	Y int

	fastdirty bool

	// Hash of the original buffer -- empty if fastdirty is on
	origHash [16]byte

	// Buffer local settings
	Settings map[string]interface{}
}

func NewBufferFromString(text, path string) *Buffer {
	return NewBuffer(strings.NewReader(text), int64(len(text)), path)
}

// NewBuffer creates a new buffer from a given reader with a given path
func NewBuffer(reader io.Reader, size int64, path string) *Buffer {
	b := new(Buffer)
	b.LineArray = NewLineArray(size, reader)

	absPath, _ := filepath.Abs(path)

	b.Path = path
	b.AbsPath = absPath

	b.Update()

	// Put the cursor at the first spot
	cursorStartY := 0
	b.Y = cursorStartY

	if size > 50000 {
		// If the file is larger than a megabyte fastdirty needs to be on
		b.fastdirty = true
	} else {
		b.origHash = md5.Sum([]byte(b.String()))
	}

	return b
}

// GetName returns buffer name
func (b *Buffer) GetName() string {
	if b.name == "" {
		if b.Path == "" {
			return "No name"
		}
		return b.Path
	}
	return b.name
}

// Update fetches the string from the rope and updates the `text` and `lines` in the buffer
func (b *Buffer) Update() {
	b.NumLines = len(b.lines)
}

// Start returns the location of the first character in the buffer
func (b *Buffer) Start() Loc {
	return Loc{0, 0}
}

// End returns the location of the last character in the buffer
func (b *Buffer) End() Loc {
	return Loc{utf8.RuneCount(b.lines[b.NumLines-1].data), b.NumLines - 1}
}

// RuneAt returns the rune at a given location in the buffer
func (b *Buffer) RuneAt(loc Loc) rune {
	line := []rune(b.Line(loc.Y))
	if len(line) > 0 {
		return line[loc.X]
	}
	return '\n'
}

// Line returns a single line
func (b *Buffer) Line(n int) string {
	if n >= len(b.lines) {
		return ""
	}
	return string(b.lines[n].data)
}

// LinesNum returns the number of lines
func (b *Buffer) LinesNum() int {
	return len(b.lines)
}

// Lines returns an array of strings containing the lines from start to end
func (b *Buffer) Lines(start, end int) []string {
	lines := b.lines[start:end]
	var slice []string
	for _, line := range lines {
		slice = append(slice, string(line.data))
	}
	return slice
}

// Len gives the length of the buffer
func (b *Buffer) Len() int {
	return Count(b.String())
}
