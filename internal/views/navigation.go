package views

import (
	"fmt"
	"lazydebrid/internal/logs"

	"github.com/jroimartin/gocui"
)

var currentViewIdx int

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
	_, err := g.SetCurrentView(ViewSearch)
	return err
}

func CycleFocusToNextView(g *gocui.Gui, v *gocui.View) error {
	currentViewIdx = (currentViewIdx + 1) % len(Views)
	name := Views[currentViewIdx]

	if _, err := g.SetCurrentView(name); err != nil {
		logs.LogEvent(fmt.Errorf("Cannot set current view to %s: %s", name, err))
		return err
	}

	err := UpdateFooter(g)
	if err != nil {
		return err
	}
	return nil
}

func CycleFocusToPreviousView(g *gocui.Gui, v *gocui.View) error {
	currentViewIdx = (currentViewIdx - 1) % len(Views)

	// wrap if we run out of views
	if currentViewIdx <= 0 {
		currentViewIdx = len(Views) - 1
	}

	name := Views[currentViewIdx]

	if _, err := g.SetCurrentView(name); err != nil {
		logs.LogEvent(fmt.Errorf("Cannot set current view to %s: %s", name, err))
		return err
	}

	err := UpdateFooter(g)
	if err != nil {
		return err
	}
	return nil
}
