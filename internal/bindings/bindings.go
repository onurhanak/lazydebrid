package bindings

import (
	"fmt"
	"lazydebrid/internal/actions"
	"lazydebrid/internal/handlers"
	"lazydebrid/internal/views"

	"github.com/jroimartin/gocui"
)

var bindErrors []error

func Keybindings(g *gocui.Gui) error {
	bind := func(viewname string, key any, mod gocui.Modifier, handler func(*gocui.Gui, *gocui.View) error) {
		if err := g.SetKeybinding(viewname, key, mod, handler); err != nil {
			bindErrors = append(bindErrors, err)
		}
	}

	bind(views.ViewActiveTorrents, 'd', gocui.ModNone, handlers.HandleDeleteTorrent)
	bind(views.ViewActiveTorrents, 's', gocui.ModNone, actions.GetTorrentStatus)
	bind(views.ViewActiveTorrents, 'j', gocui.ModNone, views.CursorDown)
	bind(views.ViewActiveTorrents, 'k', gocui.ModNone, views.CursorUp)
	bind(views.ViewActiveTorrents, gocui.KeyArrowDown, gocui.ModNone, views.CursorDown)
	bind(views.ViewActiveTorrents, gocui.KeyArrowUp, gocui.ModNone, views.CursorUp)
	bind(views.ViewActiveTorrents, '?', gocui.ModNone, handlers.ShowHelpModal)

	bind(views.ViewDetails, 'j', gocui.ModNone, views.CursorDown)
	bind(views.ViewDetails, 'k', gocui.ModNone, views.CursorUp)
	bind(views.ViewDetails, gocui.KeyArrowDown, gocui.ModNone, views.CursorDown)
	bind(views.ViewDetails, gocui.KeyArrowUp, gocui.ModNone, views.CursorUp)
	bind(views.ViewDetails, gocui.KeyEnter, gocui.ModNone, handlers.HandleDownloadSelectedFile)
	bind(views.ViewDetails, 'd', gocui.ModNone, handlers.HandleDownloadSelectedFile)
	bind(views.ViewDetails, 'D', gocui.ModNone, handlers.HandleDownloadAll)
	bind(views.ViewDetails, 'y', gocui.ModNone, views.CopyDownloadLink)
	bind(views.ViewDetails, '/', gocui.ModNone, views.FocusSearchBar)
	bind(views.ViewDetails, '?', gocui.ModNone, handlers.ShowHelpModal)

	bind(views.ViewTorrents, 'j', gocui.ModNone, views.CursorDown)
	bind(views.ViewTorrents, 'k', gocui.ModNone, views.CursorUp)
	bind(views.ViewTorrents, gocui.KeyArrowDown, gocui.ModNone, views.CursorDown)
	bind(views.ViewTorrents, gocui.KeyArrowUp, gocui.ModNone, views.CursorUp)
	bind(views.ViewTorrents, gocui.KeyEnter, gocui.ModNone, handlers.HandleTorrentFileContents)
	bind(views.ViewTorrents, '/', gocui.ModNone, views.FocusSearchBar)
	bind(views.ViewTorrents, '?', gocui.ModNone, handlers.ShowHelpModal)
	bind(views.ViewTorrents, 'd', gocui.ModNone, handlers.HandleDeleteTorrent)

	bind(views.ViewSearch, gocui.KeyEnter, gocui.ModNone, views.SearchKeyPress)

	bind(views.ViewInfo, '?', gocui.ModNone, handlers.ShowHelpModal)

	bind("", gocui.KeyCtrlA, gocui.ModNone, handlers.ShowAddMagnetModal)
	bind("", gocui.KeyCtrlP, gocui.ModNone, handlers.ShowSetPathModal)
	bind("", gocui.KeyCtrlX, gocui.ModNone, handlers.ShowSetTokenModal)
	bind("", gocui.KeyCtrlQ, gocui.ModNone, handlers.Quit)
	bind("", gocui.KeyTab, gocui.ModNone, views.CycleFocusToNextView)
	// this causes problems and should not work in search bar
	bind("", gocui.KeyBackspace2, gocui.ModNone, views.CycleFocusToPreviousView)

	if len(bindErrors) > 0 {
		return fmt.Errorf("some keybindings failed: %v", bindErrors)
	}

	return nil
}
