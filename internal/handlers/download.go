package handlers

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/jroimartin/gocui"

	"lazydebrid/internal/actions"
	"lazydebrid/internal/config"
	"lazydebrid/internal/logs"
	"lazydebrid/internal/logui"
	"lazydebrid/internal/models"
	"lazydebrid/internal/views"
)

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

	infoView := views.GetView(g, views.ViewInfo)
	now := logs.GetNow()

	logui.LogInfo(infoView, now, fmt.Sprintf("Downloading %s to %s", downloadItem.Filename, config.DownloadPath()))

	go func(item models.Download) {
		if actions.DownloadFile(item) {
			now := logs.GetNow()
			logui.LogInfo(infoView, now, fmt.Sprintf("Downloaded %s to %s", item.Filename, config.DownloadPath()))
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
		logui.LogError(v, logs.GetNow(), "Failed to copy download link", err)
	}
	return nil
}
