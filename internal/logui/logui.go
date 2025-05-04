package logui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func LogInfo(v *gocui.View, time string, errorString string) {
	fmt.Fprintf(v, "[ %s ] %s", time, errorString)
}

func LogError(v *gocui.View, time string, errorString string, err error) {
	fmt.Fprintf(v, "[ %s ] %s\n%s", time, errorString, err)
}
