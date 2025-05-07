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

func handleStartDownload(g *gocui.Gui, torrentFile models.TorrentFileDetailed) {
	log := func(msg string, success bool, err error) {
		g.Update(func(g *gocui.Gui) error {
			views.UpdateUILog(g, msg, success, err)
			return nil
		})
	}

	log(fmt.Sprintf("Downloading %s to %s", torrentFile.Filename, config.DownloadPath()), true, nil)

	if actions.DownloadFile(torrentFile) {
		log(fmt.Sprintf("Downloaded %s to %s", torrentFile.Filename, config.DownloadPath()), true, nil)
	} else {
		log(fmt.Sprintf("Failed to download %s", torrentFile.Filename), false, fmt.Errorf("download failed"))
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
