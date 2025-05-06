package views

import (
	"fmt"
	"lazydebrid/internal/config"
	"lazydebrid/internal/data"
	"lazydebrid/internal/logs"
	"lazydebrid/internal/models"
	"lazydebrid/internal/utils"
	"log"
	"strings"

	"github.com/atotto/clipboard"
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
		log.Println(fmt.Errorf("Cannot delete view %s: %s", name, err))
	}
	_, err := g.SetCurrentView(ViewTorrents)
	if err != nil {

		logs.LogEvent(fmt.Errorf("Cannot set current view to %s: %s", name, err))
		logs.LogEvent(err)
	}
	return err
}

func GetSelectedTorrentID(v *gocui.View) (string, error) {
	_, cy := v.Cursor()
	if cy < 0 {
		return "", fmt.Errorf("cursor is off-screen or uninitialized")
	}
	if cy >= len(data.TorrentLineIndex) {
		return "", fmt.Errorf("cursor index %d out of bounds (max %d)", cy, len(data.TorrentLineIndex)-1)
	}
	return data.TorrentLineIndex[cy], nil
}

func GetSelectedLine(v *gocui.View) (string, error) {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return "", fmt.Errorf("unable to get selected line: %w", err)
	}
	return line, nil
}

func GetSelectedItem(v *gocui.View) (models.Download, error) {
	line, err := GetSelectedLine(v)
	if err != nil {
		return models.Download{}, fmt.Errorf("unable to get selected line: %w", err)
	}
	item, ok := data.FilesMap[line]
	if !ok {
		return models.Download{}, fmt.Errorf("no download item found for selected line")
	}
	return item, nil
}

func PopulateViews(g *gocui.Gui) {
	torrentsView := GetView(g, ViewTorrents)
	torrentsView.Clear()
	for _, item := range data.UserDownloads {
		if item.Status == "downloaded" {
			fmt.Fprintln(torrentsView, strings.TrimSpace(item.Filename))
		}
	}

	activeView := GetView(g, ViewActiveTorrents)
	activeView.Clear()

	// Add active downloads from the API
	for _, item := range data.UserDownloads {
		if item.Status == "queued" || item.Status == "downloading" {
			fmt.Fprintln(activeView, item.ID)
		}
	}
	// Add active downloads from the present session
	for _, item := range data.ActiveDownloads {
		fmt.Fprintln(activeView, item.ID)
	}

	UpdateFooter(g)
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
	mainView.Highlight = false

	torrentItem, ok := data.DownloadMap[strings.TrimSpace(line)]
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

func CopyDownloadLink(g *gocui.Gui, v *gocui.View) error {
	item, err := GetSelectedItem(v)
	if err != nil {
		return err
	}

	if err := clipboard.WriteAll(item.Download); err != nil {
		UpdateUILog(g, fmt.Sprintf("Failed to copy download link: %s", err), false, err)
		return err
	}

	UpdateUILog(g, fmt.Sprintf("Copied download link for %s", item.Filename), true, nil)
	return nil
}

func SearchKeyPress(g *gocui.Gui, v *gocui.View) error {
	config.SetSearchQuery(v.Buffer())

	if err := utils.RenderList(g); err != nil {
		return err
	}

	torrentsView, _ := g.View(ViewTorrents)
	if err := UpdateDetails(g, torrentsView); err != nil {
		return err
	}

	_, err := g.SetCurrentView(ViewTorrents)
	return err
}

func ShowTorrentFiles(g *gocui.Gui, v *gocui.View, fileMap map[string]models.Download) {

	g.Update(func(g *gocui.Gui) error {
		detailsView := GetView(g, ViewDetails)
		if detailsView == nil {
			err := fmt.Errorf("torrentsView is nil")
			logs.LogEvent(err)
			return err
		}

		detailsView.Clear()
		detailsView.Highlight = true
		for key := range fileMap {
			fmt.Fprintln(detailsView, strings.TrimSpace(key))
		}

		_, _ = g.SetCurrentView(ViewDetails)
		return nil
	})
}
func DeleteCurrentView(g *gocui.Gui, v *gocui.View) error {
	currentView := g.CurrentView()
	if currentView == nil {
		return nil
	}

	switch currentView.Name() {
	case ViewTorrents, ViewDetails, ViewInfo, ViewFooter, ViewActiveTorrents, ViewSearch:
		return nil
	default:
		if err := g.DeleteView(currentView.Name()); err != nil {
			return err
		}
		_, err := g.SetCurrentView(ViewTorrents)
		return err
	}
}
