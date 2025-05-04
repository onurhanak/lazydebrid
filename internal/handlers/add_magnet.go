package handlers

import (
	"fmt"
	"strings"

	"lazydebrid/internal/actions"
	"lazydebrid/internal/logs"
	"lazydebrid/internal/logui"
	"lazydebrid/internal/views"

	"github.com/jroimartin/gocui"
)

func HandleAddMagnetLink(g *gocui.Gui, input string) error {
	input = strings.TrimSpace(input)
	info := views.GetView(g, views.ViewInfo)
	now := logs.GetNow()
	if input == "" {
		logui.LogInfo(info, now, "Error: Empty magnet link")
		return nil
	}

	id, err := actions.SendLinkToAPI(input)
	if err != nil {
		logui.LogError(info, now, "", err)
		return views.CloseView(g, views.ViewAddMagnet)
	}
	logui.LogInfo(info, now, fmt.Sprintf("Magnet added: %s", id))

	if actions.AddFilesToDebrid(id) {
		logui.LogInfo(info, now, fmt.Sprintf("All files selected for download: %s", id))

	} else {
		logui.LogError(info, now, fmt.Sprintf("Failed to select files for %s", id), nil)
	}
	return views.CloseView(g, views.ViewAddMagnet)
}
