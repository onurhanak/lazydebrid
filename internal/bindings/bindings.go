package bindings

import (
	"log"

	"lazydebrid/internal/actions"
	"lazydebrid/internal/handlers"

	"github.com/jroimartin/gocui"
)

func Keybindings(g *gocui.Gui) error {
	bind := func(viewname string, key interface{}, mod gocui.Modifier, handler func(*gocui.Gui, *gocui.View) error) {
		if err := g.SetKeybinding(viewname, key, mod, handler); err != nil {
			log.Fatalf("binding failed: %v", err)
		}
	}

	bind("activeTorrents", 'd', gocui.ModNone, actions.DeleteTorrent)
	bind("activeTorrents", 's', gocui.ModNone, actions.GetTorrentStatus)
	bind("activeTorrents", 'j', gocui.ModNone, handlers.CursorDown)
	bind("activeTorrents", 'k', gocui.ModNone, handlers.CursorUp)
	bind("torrents", 'j', gocui.ModNone, handlers.CursorDown)
	bind("torrents", 'k', gocui.ModNone, handlers.CursorUp)
	bind("", gocui.KeyArrowDown, gocui.ModNone, handlers.CursorDown)
	bind("", gocui.KeyArrowUp, gocui.ModNone, handlers.CursorUp)
	bind("torrents", gocui.KeyEnter, gocui.ModNone, handlers.DownloadSelected)
	bind("torrents", '/', gocui.ModNone, handlers.FocusSearchBar)
	bind("details", '/', gocui.ModNone, handlers.FocusSearchBar)
	bind("search", gocui.KeyEnter, gocui.ModNone, handlers.SearchKeyPress)
	bind("", gocui.KeyCtrlC, gocui.ModNone, handlers.CopyDownloadLink)
	bind("", gocui.KeyCtrlD, gocui.ModNone, handlers.DownloadSelected)
	bind("", gocui.KeyCtrlA, gocui.ModNone, handlers.ShowAddMagnetModal)
	bind("", gocui.KeyCtrlP, gocui.ModNone, handlers.ShowSetPathModal)
	bind("", gocui.KeyCtrlX, gocui.ModNone, handlers.ShowSetTokenModal)
	bind("", gocui.KeyCtrlQ, gocui.ModNone, handlers.Quit)
	bind("", gocui.KeyTab, gocui.ModNone, handlers.CycleViewHandler)
	bind("torrents", '?', gocui.ModNone, handlers.ShowHelpModal)
	return nil
}
