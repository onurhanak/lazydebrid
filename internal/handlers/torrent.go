package handlers

import (
	"fmt"
	"lazydebrid/internal/actions"
	"lazydebrid/internal/models"
	"lazydebrid/internal/views"
	"strings"

	"github.com/jroimartin/gocui"
)

func HandleDeleteTorrent(g *gocui.Gui, v *gocui.View) error {

	// check for view to handle getting torrent id differently
	currentViewName := g.CurrentView().Name()
	var torrentID string
	var err error
	var cy int
	var torrent models.Torrent
	switch currentViewName {
	case views.ViewActiveTorrents:
		torrentID, err = views.GetSelectedActiveDownload(v)
		if err != nil || strings.TrimSpace(torrentID) == "" {
			return fmt.Errorf("no torrent selected")
		}
	case views.ViewTorrents:
		torrent, cy, err = views.GetSelectedTorrent(v)

		if err != nil || strings.TrimSpace(torrent.ID) == "" {
			return fmt.Errorf("no torrent selected")
		}
		torrentID = torrent.ID

	}

	err = actions.DeleteTorrent(torrentID, cy, currentViewName)
	if err != nil {
		views.UpdateUILog(g, fmt.Sprintf("Failed to delete torrent: %s\nError: %s", torrentID, err), false, nil)
		return err
	}

	g.Update(func(g *gocui.Gui) error {

		views.UpdateUILog(g, fmt.Sprintf("Deleted torrent: %s", torrentID), true, nil)
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
