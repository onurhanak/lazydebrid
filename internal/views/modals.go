package views

import (
	"fmt"
	"lazydebrid/internal/logs"
	"strings"

	"github.com/jroimartin/gocui"
)

type InputHandler func(string) error

func ShowModal(g *gocui.Gui, name, title, content string, onSubmit InputHandler) error {
	maxX, maxY := g.Size()

	// adjust dimension based on content
	lines := strings.Count(content, "\n") + 2
	maxContentWidth := 60
	w := maxContentWidth + 4
	h := lines + 4

	x0 := (maxX - w) / 2
	y0 := (maxY - h) / 2
	x1 := x0 + w
	y1 := y0 + h

	v, err := g.SetView(name, x0, y0, x1, y1)
	if err != nil && err != gocui.ErrUnknownView {
		logs.LogEvent(fmt.Errorf("Cannot set view to %s: %s", name, err))
		return err
	}

	v.Title = title
	v.Wrap = true
	v.Editable = (name != ViewHelp)

	if name == ViewHelp {
		v.Clear()
		_, _ = fmt.Fprint(v, content)
	}

	g.DeleteKeybindings(name)
	_ = g.SetKeybinding(name, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		input := strings.TrimSpace(v.Buffer())
		_ = g.DeleteView(name)
		_, _ = g.SetCurrentView(ViewTorrents)
		return onSubmit(input)
	})

	_ = g.SetKeybinding(name, gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, _ = g.SetCurrentView(ViewTorrents)
		_ = g.DeleteView(name)
		return nil
	})

	_, err = g.SetCurrentView(name)
	g.Update(func(g *gocui.Gui) error {
		UpdateFooter(g)
		return nil
	})
	return err
}
