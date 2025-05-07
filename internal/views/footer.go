package views

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

var footerKeyMap = map[string]string{
	ViewTorrents:       TorrentsKeys,
	ViewSearch:         SearchKeys,
	ViewActiveTorrents: ActiveDownloadsKeys,
	ViewDetails:        DetailsKeys,
	ViewAddMagnet:      ModalKeys,
	ViewHelp:           ModalKeys,
}

func UpdateFooter(g *gocui.Gui) error {
	current := g.CurrentView()
	if current == nil {
		return fmt.Errorf("no current view")
	}

	viewName := current.Name()

	keysView, err := g.View(ViewFooter)
	if err != nil {
		return fmt.Errorf("cannot get footer view: %w", err)
	}
	keysView.Clear()

	if keys, ok := footerKeyMap[viewName]; ok {
		fmt.Fprint(keysView, keys)
	} else {
		fmt.Fprint(keysView, MainKeys)
	}

	return nil
}
