package handlers

import (
	"fmt"
	"strings"

	"lazydebrid/internal/actions"
	"lazydebrid/internal/views"

	"github.com/jroimartin/gocui"
)

func HandleAddMagnetLink(g *gocui.Gui, input string) error {
	input = strings.TrimSpace(input)
	if input == "" {
		views.UpdateUILog(g, "Error: Empty magnet link", true, nil)
		return nil
	}

	id, err := actions.SendLinkToAPI(input)
	if err != nil {
		views.UpdateUILog(g, "", false, err)
		return err
	}
	views.UpdateUILog(g, fmt.Sprintf("Magnet added: %s", id), true, nil)

	if actions.AddFilesToDebrid(id) {
		views.UpdateUILog(g,
			fmt.Sprintf("All files selected for download: %s", id),
			true,
			nil)
		// update activeTorrentsView
		views.PopulateViews(g)

	} else {

		views.UpdateUILog(g,
			fmt.Sprintf("Failed to select files for %s", id),
			false,
			nil)
	}
	return nil
}
