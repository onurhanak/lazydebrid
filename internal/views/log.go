package views

import (
	"fmt"
	"lazydebrid/internal/logs"

	"github.com/jroimartin/gocui"
)

func LogInfo(v *gocui.View, errorString string) {
	time := logs.GetNow()
	fmt.Fprintf(v, "\n[ %s ] %s", time, errorString)
}

func LogError(v *gocui.View, errorString string, err error) {
	time := logs.GetNow()
	fmt.Fprintf(v, "\n[ %s ]\n%s %s", time, errorString, err)
}

func UpdateUILog(g *gocui.Gui, message string, isInfo bool, err error) {
	infoView := GetView(g, ViewInfo)
	g.Update(func(g *gocui.Gui) error {
		if isInfo {
			LogInfo(infoView, message)
		} else {
			LogError(infoView, message, err)
		}
		return nil
	})
}
