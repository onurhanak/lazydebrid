package handlers

import (
	"lazydebrid/internal/actions"

	"github.com/jroimartin/gocui"
)

func FileContentsHandler(g *gocui.Gui, v *gocui.View) error {
	torrentFiles := actions.GetTorrentContents(g, v)
	RefreshTorrentsView(g, v, torrentFiles)
	return nil
}
