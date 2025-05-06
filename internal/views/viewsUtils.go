package views

import (
	"fmt"
	"lazydebrid/internal/data"
	"lazydebrid/internal/logs"
	"lazydebrid/internal/models"
	"lazydebrid/internal/utils"
	"log"
	"strings"

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

func GetSelectedItem(v *gocui.View) (models.Download, error) {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
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
