package logui

import (
	"fmt"
	"lazydebrid/internal/logs"
	"lazydebrid/internal/views"

	"github.com/jroimartin/gocui"
)

func LogInfo(v *gocui.View, time string, errorString string) {
	fmt.Fprintf(v, "\n[ %s ] %s", time, errorString)
}

func LogError(v *gocui.View, time string, errorString string, err error) {
	fmt.Fprintf(v, "\n[ %s ]\n%s %s", time, errorString, err)
}

func UpdateUILog(g *gocui.Gui, message string) {
	g.Update(func(g *gocui.Gui) error {
		infoView := views.GetView(g, views.ViewInfo)
		now := logs.GetNow()
		LogInfo(infoView, now, message)
		return nil
	})
}
