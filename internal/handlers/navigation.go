package handlers

import (
	"lazydebrid/internal/views"

	"github.com/jroimartin/gocui"
)

func CursorDown(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy+1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	currentView := g.CurrentView()
	if currentView.Name() == "torrents" {
		return UpdateDetails(g, v)
	}
	return nil
}

func CursorUp(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	if cy > 0 {
		if err := v.SetCursor(cx, cy-1); err != nil {
			return err
		}
	}

	currentView := g.CurrentView()
	if currentView.Name() == "torrents" {

		return UpdateDetails(g, v)
	}
	return nil
}

func FocusSearchBar(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView(views.ViewSearch)
	return err
}

func CycleViewHandler(g *gocui.Gui, v *gocui.View) error {
	return views.CycleFocusToNextView(g)
}
