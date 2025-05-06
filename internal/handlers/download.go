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

func handleStartDownload(g *gocui.Gui, item models.Download) {
	log := func(msg string, success bool, err error) {
		views.UpdateUILog(g, msg, success, err)
	}

	log(fmt.Sprintf("Downloading %s to %s", item.Filename, config.DownloadPath()), true, nil)

	if actions.DownloadFile(item) {
		log(fmt.Sprintf("Downloaded %s to %s", item.Filename, config.DownloadPath()), true, nil)
	} else {
		log(fmt.Sprintf("Failed to download %s", item.Filename), false, fmt.Errorf("download failed"))
	}
}

func HandleDownloadAll(g *gocui.Gui, _ *gocui.View) error {

	for _, item := range data.FilesMap {
		go func(dlItem models.Download) {
			handleStartDownload(g, dlItem)
		}(item)
	}

	return nil
}

func HandleDownloadSelectedFile(g *gocui.Gui, v *gocui.View) error {
	item, err := views.GetSelectedItem(v)
	if err != nil {
		return err
	}
	go handleStartDownload(g, item)
	return nil
}
