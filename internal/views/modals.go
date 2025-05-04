package views

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

type InputHandler func(string) error

func ShowModal(g *gocui.Gui, name, title, content string, onSubmit InputHandler) error {
	maxX, maxY := g.Size()
	w, h := maxX/2, 5
	x0 := (maxX - w) / 2
	y0 := (maxY - h) / 2
	x1 := x0 + w
	y1 := y0 + h

	v, err := g.SetView(name, x0, y0, x1, y1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	if err == nil {
		v.Title = title
		v.Wrap = true
		v.Editable = (name != ViewHelp)

		if name == ViewHelp {
			_, _ = fmt.Fprint(v, content)
		}

		g.DeleteKeybindings(name)
		g.SetKeybinding(name, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			input := strings.TrimSpace(v.Buffer())
			_ = g.DeleteView(name)
			_, _ = g.SetCurrentView(ViewTorrents)
			return onSubmit(input)
		})
	}

	_, err = g.SetCurrentView(name)
	return err
}
