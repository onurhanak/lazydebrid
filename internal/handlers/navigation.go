package handlers

import (
	"lazydebrid/internal/views"

	"github.com/jroimartin/gocui"
)

func CursorDown(g *gocui.Gui, v *gocui.View) error {
	numLines := len(v.BufferLines())
	cx, cy := v.Cursor()
	ox, oy := v.Origin()

	cursorLine := cy + oy

	// do not go down if we are at bottom
	if cursorLine >= numLines-1 {
		return nil
	}

	if err := v.SetCursor(cx, cy+1); err != nil {
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}

	if g.CurrentView().Name() == "torrents" {
		return UpdateDetails(g, v)
	}
	return nil
}

func CursorUp(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()

	if cy > 0 {
		if err := v.SetCursor(cx, cy-1); err != nil {
			return err
		}
	} else if oy > 0 { // keep going up if there are more items above
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}

	if g.CurrentView().Name() == "torrents" {
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
