package ui

import (
	"sync"

	"github.com/gdamore/tcell/views"
)

type Panel struct {
	sb   *StatusBar
	once sync.Once
	app  *App

	views.Panel
}

func (p *Panel) SetTitle(title string) {
	p.sb.SetCenter(title)
}

func (p *Panel) Init(app *App) {
	p.once.Do(func() {
		p.app = app

		p.sb = NewStatusBar()
		p.sb.SetMode(" NORMAL ")

		p.Panel.SetStatus(p.sb)
		// p.Panel.SetMenu(p.sb)
		// p.Panel.SetStatus(p.kb)
	})
}

func (p *Panel) App() *App {
	return p.app
}

func NewPanel(app *App) *Panel {
	p := &Panel{}
	p.Init(app)
	return p
}
