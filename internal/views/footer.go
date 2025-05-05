package views

import (
	"fmt"
	"lazydebrid/internal/logs"

	"github.com/jroimartin/gocui"
)

func UpdateFooter(g *gocui.Gui) error {
	current := g.CurrentView()
	if current == nil {
		return fmt.Errorf("no current view")
	}

	viewName := current.Name()
	keysView, err := g.View(ViewFooter)
	if err != nil {
		logs.LogEvent(fmt.Errorf("Cannot get keysView: %s", err))
		return err
	}
	keysView.Clear()

	switch viewName {
	case ViewTorrents:
		fmt.Fprint(keysView, TorrentsKeys)
	case ViewSearch:
		fmt.Fprint(keysView, SearchKeys)
	case ViewActiveTorrents:
		fmt.Fprint(keysView, ActiveDownloadsKeys)
	case ViewDetails:
		fmt.Fprint(keysView, DetailsKeys)
	case ViewAddMagnet, ViewHelp:
		fmt.Fprint(keysView, ModalKeys)
	default:
		fmt.Fprint(keysView, MainKeys)
	}
	return nil
}
