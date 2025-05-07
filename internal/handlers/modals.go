package handlers

import (
	"fmt"

	"lazydebrid/internal/config"
	"lazydebrid/internal/views"

	"github.com/jroimartin/gocui"
)

func ShowSetPathModal(g *gocui.Gui, v *gocui.View) error {
	return views.ShowModal(g, views.ViewSetPath, "Set Download Path", "", func(input string) error {
		if err := config.SaveSetting("downloadPath", input); err != nil {
			return err
		}
		g.Update(func(g *gocui.Gui) error {
			views.UpdateUILog(g, "Download path updated.", nil)
			return nil
		})
		return nil
	})
}

func ShowSetTokenModal(g *gocui.Gui, v *gocui.View) error {
	return views.ShowModal(g, views.ViewSetToken, "Set API Token", "", func(input string) error {
		if err := config.SaveSetting("apiToken", input); err != nil {
			return fmt.Errorf("failed to save API token: %w", err)
		}
		g.Update(func(g *gocui.Gui) error {
			views.UpdateUILog(g, "API token updated.", nil)
			return nil
		})
		return nil
	})
}

func ShowAddMagnetModal(g *gocui.Gui, v *gocui.View) error {
	return views.ShowModal(g, views.ViewAddMagnet, "Add Magnet Link", "", func(input string) error {
		return HandleAddMagnetLink(g, input)
	})
}

func ShowHelpModal(g *gocui.Gui, v *gocui.View) error {
	content := `
  ── Navigation ─────────────
  ↑ ↓, j k       Move cursor
  TAB            Pane forward
  BACKSPACE      Pane back
  /              Focus search

  ── Actions ────────────────
  ENTER, d       Download selected (Details)
  D              Download all files (Details)
  ENTER          View torrent files (Torrents)
  y              Copy download link (Details)
  a              Add magnet link (Torrents)

  ── Management ─────────────
  d              Delete torrent (Torrents / Active)
  s              Check status (Active)
  ^A             Add magnet modal (global)
  ^P             Set download path
  ^X             Set API key
  ^Q             Quit application
  ^C             Close modal

  ── Help ───────────────────
  ?              Show this help
`
	return views.ShowModal(g, views.ViewHelp, "Shortcuts", content, func(string) error { return nil })
}
