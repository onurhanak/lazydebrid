package handlers

import (
	"lazydebrid/internal/logs"
	"lazydebrid/internal/logui"
	"lazydebrid/internal/views"

	"github.com/jroimartin/gocui"
)

func UpdateUILog(g *gocui.Gui, message string) {
	g.Update(func(g *gocui.Gui) error {
		infoView := views.GetView(g, views.ViewInfo)
		now := logs.GetNow()
		logui.LogInfo(infoView, now, message)
		return nil
	})
}
