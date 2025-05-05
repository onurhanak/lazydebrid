package handlers

import (
	"fmt"
	"lazydebrid/internal/actions"
	"lazydebrid/internal/logs"
	"lazydebrid/internal/logui"
	"lazydebrid/internal/models"
	"lazydebrid/internal/views"
	"strings"

	"github.com/jroimartin/gocui"
)

func showTorrentFiles(g *gocui.Gui, v *gocui.View, fileMap map[string]models.Download) {

	g.Update(func(g *gocui.Gui) error {
		detailsView := views.GetView(g, views.ViewDetails)
		if detailsView == nil {
			err := fmt.Errorf("torrentsView is nil")
			logs.LogEvent(err)
			return err
		}

		detailsView.Clear()
		detailsView.Highlight = true
		for key := range fileMap {
			fmt.Fprintln(detailsView, strings.TrimSpace(key))
		}

		_, _ = g.SetCurrentView(views.ViewDetails)
		return nil
	})
}
func FileContentsHandler(g *gocui.Gui, v *gocui.View) error {
	logui.UpdateUILog(g, "Getting file contents...", true, nil)

	go func() {
		torrentFiles := actions.GetTorrentContents(g, v)

		g.Update(func(g *gocui.Gui) error {
			showTorrentFiles(g, v, torrentFiles)
			return nil
		})
	}()
	return nil
}
