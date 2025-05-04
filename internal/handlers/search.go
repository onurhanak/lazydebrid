package handlers

import (
	"github.com/jroimartin/gocui"

	"lazydebrid/internal/config"
	"lazydebrid/internal/utils"
	"lazydebrid/internal/views"
)

func SearchKeyPress(g *gocui.Gui, v *gocui.View) error {
	config.SetSearchQuery(v.Buffer())

	if err := utils.RenderList(g); err != nil {
		return err
	}

	torrentsView, _ := g.View(views.ViewTorrents)
	if err := UpdateDetails(g, torrentsView); err != nil {
		return err
	}

	_, err := g.SetCurrentView(views.ViewTorrents)
	return err
}
