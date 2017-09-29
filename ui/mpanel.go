package ui

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

var (
	StyleNormal = tcell.StyleDefault.
			Foreground(tcell.ColorSilver).
			Background(tcell.ColorBlack)
	StyleGood = tcell.StyleDefault.
			Foreground(tcell.ColorGreen).
			Background(tcell.ColorBlack)
	StyleWarn = tcell.StyleDefault.
			Foreground(tcell.ColorYellow).
			Background(tcell.ColorBlack)
	StyleError = tcell.StyleDefault.
			Foreground(tcell.ColorMaroon).
			Background(tcell.ColorBlack)
)

type MainPanel struct {
	name    string // service name
	err     error  // last error retrieving state
	content *views.CellView
	lines   []string
	styles  []tcell.Style
	curx    int
	cury    int

	Panel
}

// mainModel provides the model for a CellArea.
type mainModel struct {
	m *MainPanel
}

func (model *mainModel) GetBounds() (int, int) {
	// This assumes that all content is displayable runes of width 1.
	m := model.m
	y := len(m.lines)
	x := 0
	for _, l := range m.lines {
		if x < len(l) {
			x = len(l)
		}
	}
	return x, y
}

func (model *mainModel) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	var ch rune
	// var style tcell.Style

	return ch, StyleNormal, nil, 1

	// m := model.m

	// if y < 0 || y >= len(m.lines) {
	// 	return ch, StyleNormal, nil, 1
	// }

	// if x >= 0 && x < len(m.lines[y]) {
	// 	ch = rune(m.lines[y][x])
	// } else {
	// 	ch = ' '
	// }
	// style = m.styles[y]
	// if m.items[y] == m.selected {
	// 	style = style.Reverse(true)
	// }
	// return ch, style, nil, 1
}

func (model *mainModel) GetCursor() (int, int, bool, bool) {
	m := model.m
	return m.curx, m.cury, true, false
}

func (model *mainModel) MoveCursor(offx, offy int) {

	// m := model.m
	// m.curx += offx
	// m.cury += offy
	// m.updateCursor(true)
}

func (model *mainModel) SetCursor(x, y int) {
	// m := model.m
	// m.curx = x
	// m.cury = y
	// m.updateCursor(true)
}

func NewMainPanel(app *App, file string) *MainPanel {
	m := &MainPanel{}

	m.Panel.Init(app)
	m.content = views.NewCellView()
	m.SetContent(m.content)

	m.content.SetModel(&mainModel{m})
	m.content.SetStyle(StyleNormal)

	m.SetTitle(file)
	// m.SetKeys([]string{"[Q] Quit"})

	return m
}
