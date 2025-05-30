package handlers

import (
	"fmt"
	"lazydebrid/internal/actions"
	"lazydebrid/internal/data"
	"lazydebrid/internal/models"
	"lazydebrid/internal/views"
	"strings"

	"github.com/jroimartin/gocui"
)

func HandleDeleteTorrent(g *gocui.Gui, v *gocui.View) error {

	currentViewName := g.CurrentView().Name()
	var torrentID string
	var err error
	var cy int
	var torrent models.Torrent

	// check for view to handle getting torrent id differently
	switch currentViewName {
	case views.ViewActiveTorrents:
		torrentID, err = views.GetSelectedActiveDownload(v)
		if err != nil || strings.TrimSpace(torrentID) == "" {
			views.UpdateUILog(g, "no torrent selected", nil)
			// return nil
		}
	case views.ViewTorrents:
		torrent, cy, err = views.GetSelectedTorrent(v)

		if err != nil || strings.TrimSpace(torrent.ID) == "" {
			views.UpdateUILog(g, "no torrent selected", nil)
			// return fmt.Errorf("no torrent selected")
		}
		torrentID = torrent.ID
	}

	err = actions.DeleteTorrent(torrentID, cy, currentViewName)
	if err != nil {
		views.UpdateUILog(g, "Failed to delete torrent:", err)
		// return err
	}

	data.UserDownloads = actions.GetUserTorrents()
	g.Update(func(g *gocui.Gui) error {

		views.UpdateUILog(g, fmt.Sprintf("Deleted torrent: %s", torrentID), nil)
		views.PopulateViews(g)
		return nil
	})
	return nil
}

func HandleTorrentFileContents(g *gocui.Gui, v *gocui.View) error {
	views.UpdateUILog(g, "Getting file contents...", nil)

	go func() {
		torrentFiles := actions.GetTorrentContents(g, v)

		g.Update(func(g *gocui.Gui) error {
			views.ShowTorrentFiles(g, v, torrentFiles)
			return nil
		})
	}()
	return nil
}
