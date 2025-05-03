package handlers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/jroimartin/gocui"

	"lazydebrid/internal/actions"
	"lazydebrid/internal/config"
	"lazydebrid/internal/models"
	"lazydebrid/internal/utils"
	"lazydebrid/internal/views"
)

var (
	currentViewIdx int
)

func HandleAddMagnetLink(g *gocui.Gui, v *gocui.View) error {
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	maxX, maxY := g.Size()
	_, err := g.SetView("addMagnet", maxX/4, maxY/4, maxX*3/4, maxY*3/4)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	addView, _ := g.View("addMagnet")
	addView.Title = "Add Magnet Link"
	addView.Wrap = true
	addView.Editable = true
	addView.Clear()
	_, err = g.SetCurrentView("addMagnet")
	if err != nil {
		return err
	}

	if err := g.SetKeybinding("addMagnet", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		magnetLink := strings.TrimSpace(v.Buffer())
		infoView, _ := g.View("info")
		now := time.Now().Format("02 Jan 2006 15:04:00")

		if magnetLink == "" {
			_, err = fmt.Fprintf(infoView, "[%s] Error: Empty magnet link\n", now)
			if err != nil {
				log.Println(err)
			}
			return nil
		}

		downloadID, err := actions.SendLinkToAPI(magnetLink)
		if err != nil {
			log.Println(err)
			_, err = fmt.Fprintf(infoView, "[%s] Failed to add magnet: %v\n", now, err)
			if err != nil {
				log.Println(err)
				return err
			}
			return err
		}
		_, err = fmt.Fprintf(infoView, "[%s] Magnet added: %s\n", now, downloadID)
		if err != nil {
			log.Println(err)
			return err
		}

		success := actions.AddFilesToDebrid(downloadID)

		if success {
			_, err = fmt.Fprintf(infoView, "[%s] All files selected for download: %s\n", now, downloadID)
			if err != nil {
				return err
			}
		} else {
			_, err := fmt.Fprintf(infoView, "[%s] Failed to select files for %s\n", now, downloadID)
			if err != nil {
				return err
			}
			log.Printf("[%s] Failed to select files for %s\n", now, downloadID)
		}

		err = g.DeleteView("addMagnet")
		if err != nil {
			return err
		}
		_, err = g.SetCurrentView("torrents")
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		if !strings.Contains(err.Error(), "duplicate") {
			return err
		}
	}

	return nil
}

func UpdateDetails(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil || strings.TrimSpace(line) == "" {
		return nil
	}

	mainView, err := g.View("details")
	if err != nil {
		return err
	}
	mainView.Clear()

	torrentItem, ok := actions.DownloadMap[strings.TrimSpace(line)]
	if !ok {
		_, err = fmt.Fprint(mainView, "No details found.")
		if err != nil {
			return err
		}
		return nil
	}

	_, err = fmt.Fprint(mainView, utils.GenerateDetailsString(torrentItem))
	if err != nil {
		return err
	}
	return nil
}

func SearchKeyPress(g *gocui.Gui, v *gocui.View) error {
	config.SearchQuery = strings.TrimSpace(v.Buffer())
	if err := utils.RenderList(g); err != nil {
		return err
	}

	torrentsView, _ := g.View("torrents")
	err := UpdateDetails(g, torrentsView)
	if err != nil {
		return err
	}

	_, err = g.SetCurrentView("torrents")
	if err != nil {
		return err
	}
	return nil
}

func CursorDown(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy+1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}

	return UpdateDetails(g, v)
}

func CursorUp(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	if cy > 0 {
		if err := v.SetCursor(cx, cy-1); err != nil {
			return err
		}
	}
	return UpdateDetails(g, v)
}

func FocusSearchBar(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetCurrentView("search")
	if err != nil {
		return err
	}
	return nil
}
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func DeleteCurrentView(g *gocui.Gui, v *gocui.View) error {
	currentView := g.CurrentView()
	log.Println(currentView.Name())
	if currentView == nil {
		return nil
	}

	if currentView.Name() != "torrents" && currentView.Name() != "details" && currentView.Name() != "info" && currentView.Name() != "footer" && currentView.Name() != "activeTorrents" && currentView.Name() != "search" {
		err := g.DeleteView(currentView.Name())
		if err != nil {
			return err
		}
		_, err = g.SetCurrentView("torrents")
		if err != nil {
			return err
		}
	}

	return nil
}

func ShowControls(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	width := 30
	height := 10

	x0 := (maxX - width) / 2
	y0 := (maxY - height) / 2
	x1 := x0 + width
	y1 := y0 + height

	if v, err := g.SetView("controls", x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Controls"
		v.Editable = false
		v.Wrap = false
		v.Frame = true
		v.Clear()
		controlsString := "TAB: Switch\n↑↓: Navigate\nENTER: Download\n/: Search\n^A: Add Magnet\n^C: Copy Link\n^P: Set Path\n^X: Set API Key\n^Q: Quit"
		_, err = fmt.Fprint(v, controlsString)
		if err != nil {
			return err
		}
		_, err = g.SetCurrentView("controls")
		if err != nil {
			return err
		}
	}

	return nil
}

func DownloadSelected(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return err
	}
	downloadItem := actions.DownloadMap[line]
	v, _ = g.View("info")
	now := time.Now().Format("02 Jan 2006 15:04:00")

	_, err = fmt.Fprintf(v, "[%s] Downloading %s to %s", now, downloadItem.Filename, config.UserDownloadPath)
	if err != nil {
		return err
	}
	go func(torrentItem models.DebridDownload) {
		if actions.DownloadFile(torrentItem) {
			_, err = fmt.Fprintf(v, "Downloaded %s to %s", torrentItem.Filename, config.UserDownloadPath)
			if err != nil {
				log.Println(err)
			}
		}
	}(downloadItem)
	return nil
}

func CopyDownloadLink(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return err
	}
	selectedItem := actions.DownloadMap[line]
	err = clipboard.WriteAll(selectedItem.Download)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func ShowSetPathModal(g *gocui.Gui, v *gocui.View) error {
	return views.ShowModal(g, "setPath", "Set Download Path", "", func(input string) {
		_ = config.SaveSetting("downloadPath", input)
	})
}

func ShowSetTokenModal(g *gocui.Gui, v *gocui.View) error {
	return views.ShowModal(g, "setToken", "Set API Token", "", func(input string) {
		_ = config.SaveSetting("apiToken", input)
	})
}

func ShowAddMagnetModal(g *gocui.Gui, v *gocui.View) error {
	return views.ShowModal(g, "addMagnet", "Add Magnet Link", "", func(input string) {
		_, _ = actions.SendLinkToAPI(input)
	})
}

func ShowHelpModal(g *gocui.Gui, v *gocui.View) error {
	content := "TAB: Switch | ↑↓: Navigate | ENTER: Download | /: Search\n^A: Add Magnet | ^C: Copy Link | ^P: Set Path\n^X: Set API Key | ^Q: Quit"

	return views.ShowModal(g, "help", "Shortcuts", content, func(string) {})
}

func NextView(g *gocui.Gui, v *gocui.View) error {
	currentViewIdx = (currentViewIdx + 1) % len(views.Views)
	name := views.Views[currentViewIdx]
	_, err := g.SetCurrentView(name)

	keysView, _ := g.View("footer")
	keysView.Clear()
	switch name {

	case "torrents":

		_, err = fmt.Fprint(keysView, views.TorrentsKeys)
		if err != nil {
			return err
		}
	case "search":

		_, err = fmt.Fprint(keysView, views.SearchKeys)
		if err != nil {
			return err
		}
	case "activeTorrents":

		_, err = fmt.Fprint(keysView, views.ActiveDownloadsKeys)
		if err != nil {
			return err
		}
	case "details":

		_, err = fmt.Fprint(keysView, "")
		if err != nil {
			return err
		}
	}

	return err
}
