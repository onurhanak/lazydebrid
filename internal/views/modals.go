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

	lines := strings.Count(content, "\n") + 2
	width := 64
	height := lines + 4

	x0 := (maxX - width) / 2
	y0 := (maxY - height) / 2
	x1 := x0 + width
	y1 := y0 + height

	g.Update(func(g *gocui.Gui) error {
		v, err := g.SetView(name, x0, y0, x1, y1)
		if err != nil && err != gocui.ErrUnknownView {
			logs.LogEvent(fmt.Errorf("Cannot set view %s: %w", name, err))
			return err
		}

		v.Clear()
		v.Title = title
		v.Wrap = true
		v.Editable = (name != ViewHelp)

		if name == ViewHelp && content != "" {
			_, _ = fmt.Fprint(v, content)
		}

		g.DeleteKeybindings(name)

		_ = g.SetKeybinding(name, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			input := strings.TrimSpace(v.Buffer())
			err := onSubmit(input)
			if err != nil {
				UpdateUILog(g, "", err)
			}
			_ = g.DeleteView(name)
			_, _ = g.SetCurrentView(ViewTorrents)
			return nil
		})

		_ = g.SetKeybinding(name, gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			_ = g.DeleteView(name)
			_, _ = g.SetCurrentView(ViewTorrents)
			return nil
		})

		_, err = g.SetCurrentView(name)
		if err != nil {
			return err
		}

		return UpdateFooter(g)
	})

	return nil
}
