package views

import (
	"fmt"
	"lazydebrid/internal/logs"

	"github.com/jroimartin/gocui"
)

func updateFooter(g *gocui.Gui, name string) error {
	keysView, err := g.View("footer")
	if err != nil {
		logs.LogEvent(err)
		return err
	}
	keysView.Clear()

	switch name {
	case ViewTorrents:
		fmt.Fprint(keysView, TorrentsKeys)
	case ViewSearch:
		fmt.Fprint(keysView, SearchKeys)
	case ViewActiveTorrents:
		fmt.Fprint(keysView, ActiveDownloadsKeys)
	case ViewDetails:
		fmt.Fprint(keysView, DetailsKeys)
	case ViewAddMagnet: //does not work?
		fmt.Fprint(keysView, ModalKeys)
	default:
		fmt.Fprint(keysView, MainKeys)
	}
	return nil
}
