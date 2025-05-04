package bindings

import (
	"log"

	"lazydebrid/internal/actions"
	"lazydebrid/internal/handlers"
	"lazydebrid/internal/views"

	"github.com/jroimartin/gocui"
)

func Keybindings(g *gocui.Gui) error {
	bind := func(viewname string, key interface{}, mod gocui.Modifier, handler func(*gocui.Gui, *gocui.View) error) {
		if err := g.SetKeybinding(viewname, key, mod, handler); err != nil {
			log.Fatalf("binding failed: %v", err)
		}
	}

	bind(views.ViewActiveTorrents, 'd', gocui.ModNone, actions.DeleteTorrent)
	bind(views.ViewActiveTorrents, 's', gocui.ModNone, actions.GetTorrentStatus)
	bind(views.ViewActiveTorrents, 'j', gocui.ModNone, handlers.CursorDown)
	bind(views.ViewActiveTorrents, 'k', gocui.ModNone, handlers.CursorUp)
	bind(views.ViewActiveTorrents, '?', gocui.ModNone, handlers.ShowHelpModal)

	bind(views.ViewDetails, 'j', gocui.ModNone, handlers.CursorDown)
	bind(views.ViewDetails, 'k', gocui.ModNone, handlers.CursorUp)
	bind(views.ViewDetails, gocui.KeyEnter, gocui.ModNone, handlers.DownloadSelectedFile)
	bind(views.ViewDetails, 'd', gocui.ModNone, handlers.DownloadSelectedFile)
	bind(views.ViewDetails, 'D', gocui.ModNone, handlers.DownloadAll)
	bind(views.ViewDetails, 'y', gocui.ModNone, handlers.CopyDownloadLink)
	bind(views.ViewDetails, '/', gocui.ModNone, handlers.FocusSearchBar)
	bind(views.ViewDetails, '?', gocui.ModNone, handlers.ShowHelpModal)

	bind(views.ViewTorrents, 'j', gocui.ModNone, handlers.CursorDown)
	bind(views.ViewTorrents, 'k', gocui.ModNone, handlers.CursorUp)
	bind(views.ViewTorrents, gocui.KeyEnter, gocui.ModNone, handlers.FileContentsHandler)
	bind(views.ViewTorrents, '/', gocui.ModNone, handlers.FocusSearchBar)
	bind(views.ViewTorrents, '?', gocui.ModNone, handlers.ShowHelpModal)

	bind(views.ViewSearch, gocui.KeyEnter, gocui.ModNone, handlers.SearchKeyPress)

	bind("", gocui.KeyCtrlA, gocui.ModNone, handlers.ShowAddMagnetModal)
	bind("", gocui.KeyCtrlP, gocui.ModNone, handlers.ShowSetPathModal)
	bind("", gocui.KeyCtrlX, gocui.ModNone, handlers.ShowSetTokenModal)
	bind("", gocui.KeyCtrlQ, gocui.ModNone, handlers.Quit)
	bind("", gocui.KeyTab, gocui.ModNone, handlers.CycleViewHandler)

	return nil
}
