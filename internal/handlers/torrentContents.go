package handlers

import (
	"lazydebrid/internal/actions"
	"lazydebrid/internal/views"

	"github.com/jroimartin/gocui"
)

func FileContentsHandler(g *gocui.Gui, v *gocui.View) error {
	views.UpdateUILog(g, "Getting file contents...", true, nil)

	go func() {
		torrentFiles := actions.GetTorrentContents(g, v)

		g.Update(func(g *gocui.Gui) error {
			views.ShowTorrentFiles(g, v, torrentFiles)
			return nil
		})
	}()
	return nil
}
