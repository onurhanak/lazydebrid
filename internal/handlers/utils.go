package handlers

import (
	"fmt"
	"lazydebrid/internal/actions"
	"lazydebrid/internal/utils"
	"lazydebrid/internal/views"
	"strings"

	"github.com/jroimartin/gocui"
)

func PopulateViews(g *gocui.Gui) {
	torrentsView := views.GetView(g, views.ViewTorrents)
	torrentsView.Clear()
	for _, item := range actions.UserDownloads {
		if item.Status == "downloaded" {
			fmt.Fprintln(torrentsView, strings.TrimSpace(item.Filename))
		}
	}

	activeView := views.GetView(g, views.ViewActiveTorrents)
	activeView.Clear()

	// Add active downloads from the API
	for _, item := range actions.UserDownloads {
		if item.Status == "queued" || item.Status == "downloading" {
			fmt.Fprintln(activeView, item.ID)
		}
	}
	// Add active downloads from the present session
	for _, item := range actions.ActiveDownloads {
		fmt.Fprintln(activeView, item.ID)
	}

	views.UpdateFooter(g, views.ViewFooter)
}

func UpdateDetails(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil || strings.TrimSpace(line) == "" {
		return nil
	}

	mainView, err := g.View("details")
	if err != nil {
		return err
	}
	mainView.Clear()
	mainView.Highlight = false

	torrentItem, ok := actions.DownloadMap[strings.TrimSpace(line)]
	if !ok {
		_, err = fmt.Fprint(mainView, "No details found.")
		if err != nil {
			return err
		}
		return nil
	}

	_, err = fmt.Fprint(mainView, utils.GenerateDetailsString(torrentItem))
	if err != nil {
		return err
	}
	return nil
}
