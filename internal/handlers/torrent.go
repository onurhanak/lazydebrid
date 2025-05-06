package handlers

import (
	"lazydebrid/internal/actions"
	"lazydebrid/internal/views"

	"github.com/jroimartin/gocui"
)

// does not delete from view
func HandleDeleteTorrent(g *gocui.Gui, v *gocui.View) error {
	actions.DeleteTorrent(g, v)

	g.Update(func(g *gocui.Gui) error {
		views.PopulateViews(g)
		return nil
	})
	return nil
}

func HandleTorrentFileContents(g *gocui.Gui, v *gocui.View) error {
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
