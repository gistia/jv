package ui

import (
	"sync"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type StatusBar struct {
	once   sync.Once
	status string
	views.SimpleStyledTextBar
}

func (k *StatusBar) Init() {
	k.once.Do(func() {
		normal := tcell.StyleDefault.
			Foreground(tcell.ColorBlack).
			Background(tcell.ColorSilver)
		alternate := tcell.StyleDefault.
			Foreground(tcell.ColorBlue).
			Background(tcell.ColorSilver).Bold(true)

		k.SimpleStyledTextBar.Init()
		k.SimpleStyledTextBar.SetStyle(normal)
		k.RegisterLeftStyle('N', normal)
		k.RegisterLeftStyle('A', alternate)
	})
}

func (k *StatusBar) SetMode(mode string) {
	k.SetLeft(mode)
}

func NewStatusBar() *StatusBar {
	sb := &StatusBar{}
	sb.Init()
	return sb
}
