package handlers

import (
	"lazydebrid/internal/config"
	"lazydebrid/internal/views"

	"github.com/jroimartin/gocui"
)

func ShowSetPathModal(g *gocui.Gui, v *gocui.View) error {
	return views.ShowModal(g, views.ViewSetPath, "Set Download Path", "", func(input string) error {
		_ = config.SaveSetting("downloadPath", input)
		return nil
	})
}

func ShowSetTokenModal(g *gocui.Gui, v *gocui.View) error {
	return views.ShowModal(g, views.ViewSetToken, "Set API Token", "", func(input string) error {
		_ = config.SaveSetting("apiToken", input)
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
