package views

import (
	"fmt"
	"lazydebrid/internal/actions"
	"strings"

	"github.com/jroimartin/gocui"
)

var (
	Views = []string{"search", "torrents", "details", "activeTorrents"}
)

func Layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()
	splitX := (maxX * 4) / 10
	infoHeight := (maxY - 3) / 4

	detailsTop := 3
	detailsBottom := detailsTop + infoHeight
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen

	activeTop := detailsBottom + 1
	activeBottom := activeTop + infoHeight

	infoTop := activeBottom + 1
	infoBottom := maxY - 4
	if v, err := g.SetView("search", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Search"
		v.Editable = true
	}

	if torrentsView, err := g.SetView("torrents", 0, 3, splitX, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		torrentsView.Title = "Downloads"
		torrentsView.Highlight = true
		torrentsView.Wrap = false
		torrentsView.SelFgColor = gocui.ColorGreen

		for _, item := range actions.UserDownloads {
			_, err := fmt.Fprintln(torrentsView, item.Filename)
			if err != nil {
				return err
			}
		}
		err = torrentsView.SetCursor(0, 0)
		if err != nil {
			return err
		}
		_, err = g.SetCurrentView("torrents")
		if err != nil {
			return err
		}
	}

	if mainView, err := g.SetView("details", splitX+1, detailsTop, maxX-1, detailsBottom); err != nil && err != gocui.ErrUnknownView {
		return err
	} else if err == nil {
		mainView.Title = "Torrent Details"
		mainView.Wrap = true

	}

	if activeTorrentsView, err := g.SetView("activeTorrents", splitX+1, activeTop, maxX-1, activeBottom); err != nil && err != gocui.ErrUnknownView {
		return err
	} else if err == nil {
		activeTorrentsView.Title = "Active Downloads"
		activeTorrentsView.Highlight = true
		activeTorrentsView.Wrap = false
		activeTorrentsView.SelFgColor = gocui.ColorGreen
		activeTorrentsView.Clear()
		for _, item := range actions.ActiveDownloads {
			_, err = fmt.Fprintln(activeTorrentsView, item.ID)
			if err != nil {
				return err
			}
		}
		err = activeTorrentsView.SetCursor(0, 0)
		if err != nil {
			return err
		}
	}

	if infoView, err := g.SetView("info", splitX+1, infoTop, maxX-1, infoBottom); err != nil && err != gocui.ErrUnknownView {
		return err
	} else if err == nil {
		infoView.Title = "Log"
		infoView.Wrap = true
		infoView.Autoscroll = true
	}

	if footerView, err := g.SetView("footer", 0, infoBottom+1, maxX-1, infoBottom+3); err != nil && err != gocui.ErrUnknownView {
		return err
	} else if err == nil {
		footerView.Frame = true
		footerView.Wrap = true
		footerView.Title = ""

		_, err = fmt.Fprint(footerView, MainKeys)
		if err != nil {
			return err
		}
	}

	return nil
}

func ShowModal(g *gocui.Gui, name, title string, content string, onSubmit func(string)) error {
	maxX, maxY := g.Size()
	w, h := maxX/2, 5
	x0 := (maxX - w) / 2
	y0 := (maxY - h) / 2
	x1 := x0 + w
	y1 := y0 + h

	if v, err := g.SetView(name, x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = title
		if name != "help" {
			v.Title = title
			v.Editable = true
			v.Wrap = true
		} else {
			v.Editable = false
			v.Wrap = false
			_, err = fmt.Fprint(v, content)
			if err != nil {
				return err
			}
		}
		g.DeleteKeybindings(name)
		_, _ = g.SetCurrentView(name)
		err = g.SetKeybinding(name, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			input := strings.TrimSpace(v.Buffer())
			err = g.DeleteView(name)
			if err != nil {
				return err
			}
			_, err = g.SetCurrentView("torrents")
			if err != nil {
				return err
			}
			onSubmit(input)
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}
