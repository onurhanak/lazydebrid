package handlers

import (
	"fmt"
	"lazydebrid/internal/actions"
	"lazydebrid/internal/logs"
	"lazydebrid/internal/models"
	"lazydebrid/internal/views"
	"log"

	"github.com/jroimartin/gocui"
)

func refreshTorrentsView(g *gocui.Gui, v *gocui.View, fileMap map[string]models.Download) {
	log.Println("Refreshing")

	g.Update(func(g *gocui.Gui) error {
		detailsView := views.GetView(g, views.ViewDetails)
		if detailsView == nil {
			err := fmt.Errorf("torrentsView is nil")
			logs.LogEvent(err)
			return err
		}

		detailsView.Clear()
		detailsView.Highlight = true
		for key, _ := range fileMap {
			fmt.Fprintln(detailsView, key)
		}

		_, _ = g.SetCurrentView(views.ViewDetails)
		return nil
	})
}

func FileContentsHandler(g *gocui.Gui, v *gocui.View) error {
	torrentFiles := actions.GetTorrentContents(g, v)
	refreshTorrentsView(g, v, torrentFiles)
	return nil
}
