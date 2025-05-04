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
		logs.LogEvent(err)
		return err
	}

	keysView, err := g.View("footer")
	if err != nil {
		logs.LogEvent(err)
		return err
	}
	keysView.Clear()

	switch name {
	case ViewTorrents:
		fmt.Fprint(keysView, TorrentsKeys)
	case ViewSearch:
		fmt.Fprint(keysView, SearchKeys)
	case ViewActiveTorrents:
		fmt.Fprint(keysView, ActiveDownloadsKeys)
	}

	return nil
}
