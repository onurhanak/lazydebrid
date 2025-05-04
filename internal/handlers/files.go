package handlers

import (
	"fmt"
	"lazydebrid/internal/models"
	"lazydebrid/internal/views"
	"log"

	"github.com/jroimartin/gocui"
)

func RefreshTorrentsView(g *gocui.Gui, v *gocui.View, fileList []models.Download) {
	log.Println("Refreshing")

	g.Update(func(g *gocui.Gui) error {
		detailsView := views.GetView(g, views.ViewDetails)
		if detailsView == nil {
			log.Println("torrentsView is nil")
			return fmt.Errorf("torrentsView is nil")
		}

		detailsView.Clear()

		for _, file := range fileList {
			fmt.Fprintln(detailsView, file.Filename)
		}

		_, _ = g.SetCurrentView(views.ViewDetails)
		return nil
	})
}
