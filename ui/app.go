package ui

import (
	"log"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

type App struct {
	app    *views.Application
	view   views.View
	panel  views.Widget
	main   *MainPanel
	logger *log.Logger

	views.WidgetWatchers
}

func (a *App) ShowMain() {
	a.show(a.main)
}

func (a *App) show(w views.Widget) {
	a.app.PostFunc(func() {
		if w != a.panel {
			a.panel.SetView(nil)
			a.panel = w
		}

		a.panel.SetView(a.view)
		a.Resize()
		a.app.Refresh()
	})
}

func NewApp(file string) *App {
	app := &App{}
	app.app = &views.Application{}
	app.main = NewMainPanel(app, file)
	app.panel = app.main

	app.app.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorSilver).
		Background(tcell.ColorBlack))

	go app.refresh()
	return app
}

func (a *App) refresh() {
}

func (a *App) Draw() {
	if a.panel != nil {
		a.panel.Draw()
	}
}

func (a *App) Resize() {
	if a.panel != nil {
		a.panel.Resize()
	}
}

func (a *App) SetView(view views.View) {
	a.view = view
	if a.panel != nil {
		a.panel.SetView(view)
	}
}

func (a *App) Quit() {
	/* This just posts the quit event. */
	a.app.Quit()
}

func (a *App) SetLogger(logger *log.Logger) {
	a.logger = logger
	if logger != nil {
		logger.Printf("Start logger")
	}
}
func (a *App) Logf(fmt string, v ...interface{}) {
	if a.logger != nil {
		a.logger.Printf(fmt, v...)
	}
}

func (a *App) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		// Intercept a few control keys up front, for global handling.
		case tcell.KeyCtrlC:
			a.Quit()
			return true
		case tcell.KeyCtrlL:
			a.app.Refresh()
			return true
		}
	}

	if a.panel != nil {
		return a.panel.HandleEvent(ev)
	}
	return false
}

func (a *App) Size() (int, int) {
	if a.panel != nil {
		return a.panel.Size()
	}
	return 0, 0
}

func (a *App) Run() {
	a.Logf("Starting up user interface")
	a.app.SetRootWidget(a)
	a.ShowMain()
	go func() {
		// Give us periodic updates
		for {
			a.app.Update()
			time.Sleep(time.Second)
		}
	}()
	a.Logf("Starting app loop")
	a.app.Run()
}
