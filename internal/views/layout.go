package views

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const (
	ViewTorrents       = "torrents"
	ViewDetails        = "details"
	ViewInfo           = "info"
	ViewSearch         = "search"
	ViewActiveTorrents = "activeTorrents"
	ViewFooter         = "footer"
	ViewAddMagnet      = "addMagnet"
	ViewSetPath        = "setPath"
	ViewSetToken       = "setToken"
	ViewHelp           = "help"
)

var (
	Views = []string{ViewSearch, ViewTorrents, ViewDetails, ViewActiveTorrents, ViewInfo}
)

var OnLayoutReady func(*gocui.Gui)

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

		// Populate later in main.go to avoid import cycle
		//
		//for _, item := range actions.UserDownloads {
		//	_, err := fmt.Fprintln(torrentsView, item.Filename)
		//	if err != nil {
		//		return err
		//	}
		//}
		err = torrentsView.SetCursor(0, 0)
		if err != nil {
			return err
		}
		_, err = g.SetCurrentView(ViewTorrents)
		updateFooter(g, ViewTorrents)
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

	if activeTorrentsView, err := g.SetView(ViewActiveTorrents, splitX+1, activeTop, maxX-1, activeBottom); err != nil && err != gocui.ErrUnknownView {
		return err
	} else if err == nil {
		activeTorrentsView.Title = "Active Downloads"
		activeTorrentsView.Highlight = true
		activeTorrentsView.Wrap = false
		activeTorrentsView.SelFgColor = gocui.ColorGreen
		activeTorrentsView.Clear()

		// Populate later in main.go to avoid import cycle
		//
		//for _, item := range actions.ActiveDownloads {
		//	_, err = fmt.Fprintln(activeTorrentsView, item.ID)
		//	if err != nil {
		//		return err
		//	}
		//}
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
		footerView.Title = "Shortcuts"

		_, err = fmt.Fprint(footerView, MainKeys)
		if err != nil {
			return err
		}
	}
	if OnLayoutReady != nil {
		OnLayoutReady(g)
	}
	return nil
}
