package views

import (
	"fmt"
	"lazydebrid/internal/logs"

	"github.com/jroimartin/gocui"
)

var currentViewIdx int

func CycleFocusToNextView(g *gocui.Gui) error {
	currentViewIdx = (currentViewIdx + 1) % len(Views)
	name := Views[currentViewIdx]

	if _, err := g.SetCurrentView(name); err != nil {
		logs.LogEvent(fmt.Errorf("Cannot set current view to %s: %s", name, err))
		return err
	}

	err := UpdateFooter(g, name)
	if err != nil {
		return err
	}
	return nil
}
