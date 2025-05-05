package handlers

import (
	"lazydebrid/internal/views"

	"github.com/jroimartin/gocui"
)

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func DeleteCurrentView(g *gocui.Gui, v *gocui.View) error {
	currentView := g.CurrentView()
	if currentView == nil {
		return nil
	}

	switch currentView.Name() {
	case views.ViewTorrents, views.ViewDetails, views.ViewInfo, views.ViewFooter, views.ViewActiveTorrents, views.ViewSearch:
		return nil
	default:
		if err := g.DeleteView(currentView.Name()); err != nil {
			return err
		}
		_, err := g.SetCurrentView(views.ViewTorrents)
		return err
	}
}
