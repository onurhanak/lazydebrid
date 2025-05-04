package handlers

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/jroimartin/gocui"

	"lazydebrid/internal/actions"
	"lazydebrid/internal/config"
	"lazydebrid/internal/logui"
	"lazydebrid/internal/models"
	"lazydebrid/internal/views"
)

func DownloadAll(g *gocui.Gui, v *gocui.View) error {

	for _, downloadItem := range actions.FilesMap {
		logui.UpdateUILog(g, views.ViewInfo, fmt.Sprintf("Downloading %s to %s", downloadItem.Filename, config.DownloadPath()))
		go func(item models.Download) {
			if actions.DownloadFile(item) {
				logui.UpdateUILog(g, views.ViewInfo, fmt.Sprintf("Downloaded %s to %s", item.Filename, config.DownloadPath()))
			}
		}(downloadItem)
	}

	return nil
}

func DownloadSelectedFile(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return err
	}

	downloadItem, ok := actions.FilesMap[line]
	if !ok {
		return fmt.Errorf("no download item found for selected line")
	}

	logui.UpdateUILog(g, views.ViewInfo, fmt.Sprintf("Downloading %s to %s", downloadItem.Filename, config.DownloadPath()))
	go func(item models.Download) {
		if actions.DownloadFile(item) {
			logui.UpdateUILog(g, views.ViewInfo, fmt.Sprintf("Downloaded %s to %s", item.Filename, config.DownloadPath()))
		}
	}(downloadItem)

	return nil
}

func CopyDownloadLink(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return err
	}

	item, ok := actions.FilesMap[line]
	if !ok {
		return fmt.Errorf("no download link found")
	}

	if err := clipboard.WriteAll(item.Download); err != nil {
		logui.UpdateUILog(g, v.Name(), fmt.Sprintf("Failed to copy download link: %s", err))
	}
	return nil
}
