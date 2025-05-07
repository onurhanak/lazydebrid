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

func ShowHelpModal(g *gocui.Gui, v *gocui.View) error {
	content := `
  ── Navigation ─────────────
  ↑ ↓       Move cursor
  TAB       Switch view
  /         Focus search

  ── Actions ────────────────
  ENTER     Download selected
  ^C        Copy download link
  D         Download all files

  ── Management ─────────────
  ^A        Add magnet link
  ^P        Set download path
  ^X        Set API key
  ^Q        Quit application
`
	return views.ShowModal(g, views.ViewHelp, "Shortcuts", content, func(string) error { return nil })
}

func ShowAddMagnetModal(g *gocui.Gui, v *gocui.View) error {
	return views.ShowModal(g, views.ViewAddMagnet, "Add Magnet Link", "", func(input string) error {
		return HandleAddMagnetLink(g, input)
	})
}
