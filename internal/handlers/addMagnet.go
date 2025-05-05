package handlers

import (
	"fmt"
	"strings"

	"lazydebrid/internal/actions"
	"lazydebrid/internal/logui"

	"github.com/jroimartin/gocui"
)

func HandleAddMagnetLink(g *gocui.Gui, input string) error {
	input = strings.TrimSpace(input)
	if input == "" {
		logui.UpdateUILog(g, "Error: Empty magnet link", true, nil)
		return nil
	}

	id, err := actions.SendLinkToAPI(input)
	if err != nil {
		logui.UpdateUILog(g, "", false, err)
		return err
	}
	logui.UpdateUILog(g, fmt.Sprintf("Magnet added: %s", id), true, nil)

	if actions.AddFilesToDebrid(id) {
		logui.UpdateUILog(g,
			fmt.Sprintf("All files selected for download: %s", id),
			true,
			nil)
		// update activeTorrentsView
		PopulateViews(g)

	} else {

		logui.UpdateUILog(g,
			fmt.Sprintf("Failed to select files for %s", id),
			false,
			nil)
	}
	return nil
}
