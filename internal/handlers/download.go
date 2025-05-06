package handlers

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/jroimartin/gocui"

	"lazydebrid/internal/actions"
	"lazydebrid/internal/config"
	"lazydebrid/internal/models"
	"lazydebrid/internal/views"
)

func startDownload(g *gocui.Gui, item models.Download) {
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

func DownloadAll(g *gocui.Gui, _ *gocui.View) error {

	for _, item := range actions.FilesMap {
		go func(dlItem models.Download) {
			startDownload(g, dlItem)
		}(item)
	}

	return nil
}

func DownloadSelectedFile(g *gocui.Gui, v *gocui.View) error {
	item, err := getSelectedItem(v)
	if err != nil {
		return err
	}
	go startDownload(g, item)
	return nil
}

func CopyDownloadLink(g *gocui.Gui, v *gocui.View) error {
	item, err := getSelectedItem(v)
	if err != nil {
		return err
	}

	if err := clipboard.WriteAll(item.Download); err != nil {
		views.UpdateUILog(g, fmt.Sprintf("Failed to copy download link: %s", err), false, err)
		return err
	}

	views.UpdateUILog(g, fmt.Sprintf("Copied download link for %s", item.Filename), true, nil)
	return nil
}
