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
		g.Update(func(g *gocui.Gui) error {
			views.UpdateUILog(g, "Error: Empty magnet link", nil)
			return nil
		})
		return nil
	}

	go func() {
		id, err := actions.SendLinkToAPI(input)
		if err != nil {
			g.Update(func(g *gocui.Gui) error {
				views.UpdateUILog(g, "Failed to add magnet:", err)
				return nil
			})
			return
		}

		g.Update(func(g *gocui.Gui) error {
			views.UpdateUILog(g, fmt.Sprintf("Magnet added: %s", id), nil)
			views.PopulateViews(g)
			return nil
		})
	}()

	return nil
}
