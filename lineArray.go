package main

import (
	"bufio"
	"io"
	"unicode/utf8"

	"github.com/Jeffail/gabs"
)

func runeToByteIndex(n int, txt []byte) int {
	if n == 0 {
		return 0
	}

	count := 0
	i := 0
	for len(txt) > 0 {
		_, size := utf8.DecodeRune(txt)

		txt = txt[size:]
		count += size
		i++

		if i == n {
			break
		}
	}
	return count
}

// LogEntry is a parsed raw line
type LogEntry struct {
	timestamp string
	level     string
	message   string
	data      map[string]interface{}
}

// Line is a raw line
type Line struct {
	data  []byte
	entry LogEntry
}

func (line *Line) String() string {
	str := " "
	str += PadRight(line.entry.timestamp, " ", 24)
	str += " "
	str += PadRight(line.entry.level, " ", 5)
	str += " "
	str += line.entry.message
	return str
}

func NewLine(data []byte) Line {
	// data := []byte(rawLine)
	// ignore errors
	parsedLine, err := gabs.ParseJSON(data)

	if err != nil {
		entry := LogEntry{"", "", "", nil}
		return Line{data, entry}
	}

	timestamp := parsedLine.Path("timestamp").Data().(string)
	level := parsedLine.Path("level").Data().(string)
	message := parsedLine.Path("message").Data().(string)
	entryData := parsedLine.Data().(map[string]interface{})
	entry := LogEntry{timestamp, level, message, entryData}
	return Line{data, entry}
}

// A LineArray simply stores and array of lines and makes it easy to insert
// and delete in it
type LineArray struct {
	lines []Line
}

func Append(slice []Line, data ...Line) []Line {
	l := len(slice)
	if l+len(data) > cap(slice) { // reallocate
		// Allocate double what's needed, for future growth.
		newSlice := make([]Line, (l+len(data))+10000)
		// The copy function is predeclared and works for any slice type.
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : l+len(data)]
	for i, c := range data {
		slice[l+i] = c
	}
	return slice
}

// NewLineArray returns a new line array from an array of bytes
func NewLineArray(size int64, reader io.Reader) *LineArray {
	la := new(LineArray)

	la.lines = make([]Line, 0, 1000)

	br := bufio.NewReader(reader)
	var loaded int

	n := 0
	for {
		data, err := br.ReadBytes('\n')
		if len(data) > 1 && data[len(data)-2] == '\r' {
			data = append(data[:len(data)-2], '\n')
			if fileformat == 0 {
				fileformat = 2
			}
		} else if len(data) > 0 {
			if fileformat == 0 {
				fileformat = 1
			}
		}

		if n >= 1000 && loaded >= 0 {
			totalLinesNum := int(float64(size) * (float64(n) / float64(loaded)))
			newSlice := make([]Line, len(la.lines), totalLinesNum+10000)
			// The copy function is predeclared and works for any slice type.
			copy(newSlice, la.lines)
			la.lines = newSlice
			loaded = -1
		}

		if loaded >= 0 {
			loaded += len(data)
		}

		if err != nil {
			if err == io.EOF {
				// la.lines = Append(la.lines, Line{data[:], nil, nil, false})
				la.lines = Append(la.lines, NewLine(data[:]))
			}
			// Last line was read
			break
		} else {
			// la.lines = Append(la.lines, Line{data[:len(data)-1], nil, nil, false})
			la.lines = Append(la.lines, NewLine(data[:len(data)-1]))
		}
		n++
	}

	return la
}

// Returns the String representation of the LineArray
func (la *LineArray) String() string {
	str := ""
	for i, l := range la.lines {
		str += string(l.data)
		if i != len(la.lines)-1 {
			str += "\n"
		}
	}
	return str
}

// SaveString returns the string that should be written to disk when
// the line array is saved
// It is the same as string but uses crlf or lf line endings depending
func (la *LineArray) SaveString(useCrlf bool) string {
	str := ""
	for i, l := range la.lines {
		str += string(l.data)
		if i != len(la.lines)-1 {
			if useCrlf {
				str += "\r"
			}
			str += "\n"
		}
	}
	return str
}

// inserts a byte at a given location
func (la *LineArray) insertByte(pos Loc, value byte) {
	la.lines[pos.Y].data = append(la.lines[pos.Y].data, 0)
	copy(la.lines[pos.Y].data[pos.X+1:], la.lines[pos.Y].data[pos.X:])
	la.lines[pos.Y].data[pos.X] = value
}

// DeleteToEnd deletes from the end of a line to the position
func (la *LineArray) DeleteToEnd(pos Loc) {
	la.lines[pos.Y].data = la.lines[pos.Y].data[:pos.X]
}

// DeleteFromStart deletes from the start of a line to the position
func (la *LineArray) DeleteFromStart(pos Loc) {
	la.lines[pos.Y].data = la.lines[pos.Y].data[pos.X+1:]
}

// DeleteLine deletes the line number
func (la *LineArray) DeleteLine(y int) {
	la.lines = la.lines[:y+copy(la.lines[y:], la.lines[y+1:])]
}

// DeleteByte deletes the byte at a position
func (la *LineArray) DeleteByte(pos Loc) {
	la.lines[pos.Y].data = la.lines[pos.Y].data[:pos.X+copy(la.lines[pos.Y].data[pos.X:], la.lines[pos.Y].data[pos.X+1:])]
}

// Substr returns the string representation between two locations
func (la *LineArray) Substr(start, end Loc) string {
	startX := runeToByteIndex(start.X, la.lines[start.Y].data)
	endX := runeToByteIndex(end.X, la.lines[end.Y].data)
	if start.Y == end.Y {
		return string(la.lines[start.Y].data[startX:endX])
	}
	var str string
	str += string(la.lines[start.Y].data[startX:]) + "\n"
	for i := start.Y + 1; i <= end.Y-1; i++ {
		str += string(la.lines[i].data) + "\n"
	}
	str += string(la.lines[end.Y].data[:endX])
	return str
}
