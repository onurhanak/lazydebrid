package views

import (
	"fmt"
	"lazydebrid/internal/logs"

	"github.com/jroimartin/gocui"
)

func GetView(g *gocui.Gui, name string) *gocui.View {
	v, _ := g.View(name)
	return v
}

func LogViewInfo(v *gocui.View, time string, errorString string) {
	fmt.Fprintf(v, "[ %s ] %s", time, errorString)
}

func LogViewError(v *gocui.View, time string, errorString string, err error) {
	fmt.Fprintf(v, "[ %s ] %s\n%s", time, errorString, err)
}

func CloseView(g *gocui.Gui, name string) error {
	if err := g.DeleteView(name); err != nil {
		return err
	}
	_, err := g.SetCurrentView(ViewTorrents)
	if err != nil {
		logs.LogEvent(err)
	}
	return err
}
