package handlers

import (
	"fmt"

	"github.com/jroimartin/gocui"

	"lazydebrid/internal/actions"
	"lazydebrid/internal/config"
	"lazydebrid/internal/data"
	"lazydebrid/internal/models"
	"lazydebrid/internal/views"
)

const maxConcurrentDownloads = 5

// TODO
// this should also check for download slot
func handleStartDownload(g *gocui.Gui, torrentFile models.TorrentFileDetailed) {
	startMsg := fmt.Sprintf("Starting download: %s → %s", torrentFile.Filename, config.DownloadPath())
	endMsg := fmt.Sprintf("Downloaded: %s", torrentFile.Filename)

	g.Update(func(g *gocui.Gui) error {
		views.UpdateUILog(g, startMsg, nil)
		return nil
	})

	if err := actions.DownloadFile(torrentFile); err == nil {
		g.Update(func(g *gocui.Gui) error {
			views.UpdateUILog(g, endMsg, nil)
			return nil
		})
	} else {
		g.Update(func(g *gocui.Gui) error {
			views.UpdateUILog(g, fmt.Sprintf("Failed to download %s", torrentFile.Filename), err)
			return nil
		})
	}
}

func HandleDownloadAll(g *gocui.Gui, _ *gocui.View) error {
	sem := make(chan struct{}, maxConcurrentDownloads)

	for _, torrentFile := range data.FilesMap {
		tf := torrentFile
		sem <- struct{}{} // block until download slot available
		go func() {
			defer func() { <-sem }()
			handleStartDownload(g, tf)
		}()
	}

	return nil
}

func HandleDownloadSelectedFile(g *gocui.Gui, v *gocui.View) error {
	torrentFile, err := views.GetSelectedTorrentFile(v)
	if err != nil {
		return err
	}
	go handleStartDownload(g, torrentFile)
	return nil
}
