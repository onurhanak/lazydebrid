package views

import (
	"fmt"
	"lazydebrid/internal/logs"

	"github.com/jroimartin/gocui"
)

func logInfo(v *gocui.View, msg string, err error) {
	time := logs.GetNow()
	if err != nil {
		fmt.Fprintf(v, "\n[ %s ]\n%s", time, err)
	} else {
		fmt.Fprintf(v, "\n[ %s ]\n%s", time, msg)
	}
}

func UpdateUILog(g *gocui.Gui, message string, err error) {
	g.Update(func(g *gocui.Gui) error {
		v := GetView(g, ViewInfo)
		if v == nil {
			return nil
		}
		logInfo(v, message, err)
		return nil
	})
}
