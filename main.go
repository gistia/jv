package main

import (
	"fmt"
	"os"

	"github.com/zyedidia/tcell"
)

var (
	// The main screen
	screen tcell.Screen

	// Object to send messages and prompts to the user
	messenger *Messenger

	view *View

	// The default highlighting style
	// This simply defines the default foreground and background colors
	defStyle tcell.Style
	events   chan tcell.Event
)

func main() {
	InitScreen()
	buffer := LoadInput()

	view = NewView(buffer)

	RedrawAll()
}

func LoadInput() *Buffer {
	filename := os.Args[1]

	var buffer *Buffer
	if _, e := os.Stat(filename); e == nil {
		input, err := os.Open(filename)
		stat, _ := input.Stat()
		defer input.Close()
		if err != nil {
			panic(err)
		}
		if stat.IsDir() {
			TermMessage("Cannot read", filename, "because it is a directory")
		}
		buffer = NewBuffer(input, FSize(input), filename)
	} else {
		TermMessage("File not found", filename)
	}

	return buffer
}

func InitScreen() {
	// Should we enable true color?
	truecolor := os.Getenv("JV_TRUECOLOR") == "1"

	// In order to enable true color, we have to set the TERM to `xterm-truecolor` when
	// initializing tcell, but after that, we can set the TERM back to whatever it was
	oldTerm := os.Getenv("TERM")
	if truecolor {
		os.Setenv("TERM", "xterm-truecolor")
	}

	// Initilize tcell
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Println(err)
		if err == tcell.ErrTermNotFound {
			fmt.Println("jv does not recognize your terminal:", oldTerm)
			fmt.Println("Please go to https://github.com/zyedidia/mkinfo to read about how to fix this problem (it should be easy to fix).")
		}
		os.Exit(1)
	}
	if err = screen.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Now we can put the TERM back to what it was before
	if truecolor {
		os.Setenv("TERM", oldTerm)
	}

	screen.SetStyle(defStyle)
}

// RedrawAll redraws everything -- all the views and the messenger
func RedrawAll() {
	messenger.Clear()

	w, h := screen.Size()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			screen.SetContent(x, y, ' ', nil, defStyle)
		}
	}

	view.Display()
	// messenger.Display()
	screen.Show()
}

// logfile := "/Users/fcoury/logs/jvg.log"
// var dlog *log.Logger
// if logfile != "" {
// 	f, e := os.Create(logfile)
// 	if e == nil {
// 		dlog = log.New(f, "DEBUG:", log.LstdFlags)
// 		log.SetOutput(f)
// 	}
// }
// if e := doUI("/Users/fcoury/logs/large.log", dlog); e != nil {
// 	fmt.Println("worked")
// } else {
// 	return
// }

// file, err := os.Open("/Users/fcoury/logs/large.log")
// defer file.Close()

// if err != nil {
// 	panic(err)
// }

// reader := bufio.NewReader(file)

// var line string
// var lines []*gabs.Container
// for {
// 	line, err = reader.ReadString('\n')
// 	if err != nil {
// 		break
// 	}
// 	jsonLine, err := gabs.ParseJSON([]byte(line))
// 	if err != nil {
// 		fmt.Printf("WARN: Could not read line '%s'", line)
// 	}
// 	lines = append(lines, jsonLine)
// }

// fmt.Printf("%v", lines)

// for line := range lines {
// 	fmt.Printf("%s %s", line.Path("timestamp").(string), line.Path("message").(string))
// }
// jsonParsed, err := gabs.ParseJSON([]byte(raw))

// children, err := jsonParsed.ChildrenMap()
// if err != nil {
// 	panic(err)
// }
// for key, child := range children {
// 	fmt.Printf("key: %v, value: %v\n", key, child.Data())
// }
