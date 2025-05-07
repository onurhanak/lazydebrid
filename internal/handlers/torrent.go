package handlers

import (
	"lazydebrid/internal/actions"
	"lazydebrid/internal/views"

	"github.com/jroimartin/gocui"
)

func HandleDeleteTorrent(g *gocui.Gui, v *gocui.View) error {
	if err := actions.DeleteTorrent(g, v); err != nil {
		return err
	}

	// does not update the view for some reason
	g.Update(func(g *gocui.Gui) error {
		views.PopulateViews(g)
		return nil
	})
	return nil
}

func HandleTorrentFileContents(g *gocui.Gui, v *gocui.View) error {
	views.UpdateUILog(g, "Getting file contents...", true, nil)

	// Make a shallow copy of view cursor or content here if needed in future
	go func() {
		torrentFiles := actions.GetTorrentContents(g, v)

		g.Update(func(g *gocui.Gui) error {
			views.ShowTorrentFiles(g, v, torrentFiles)
			return nil
		})
	}()
	return nil
}
