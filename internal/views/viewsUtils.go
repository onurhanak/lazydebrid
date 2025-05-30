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

const (
	ColorReset  = "\033[0m"
	ColorCyan   = "\033[36m"
	ColorYellow = "\033[33m"
	ColorGreen  = "\033[32m"
	ColorBlue   = "\033[34m"
)

// without this, if user deletes all downloaded torrents
// the ui log will be spammed infinitely
// until a download is added
var warnedNoTorrents bool

func GenerateDetailsString(torrentItem models.Torrent) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%sID        :%s %s\n", ColorYellow, ColorReset, torrentItem.ID))
	sb.WriteString(fmt.Sprintf("%sFilename  :%s %s\n", ColorYellow, ColorReset, torrentItem.Filename))
	sb.WriteString(fmt.Sprintf("%sFilesize  :%s %s\n", ColorYellow, ColorReset, humanReadableBytes(torrentItem.Bytes)))

	return sb.String()
}

func humanReadableBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func RenderList(g *gocui.Gui) error {
	v, err := g.View("torrents")
	if err != nil {
		return err
	}
	v.Clear()
	err = v.SetCursor(0, 0)
	if err != nil {
		return err
	}

	for _, torrentItem := range data.UserDownloads {
		if config.SearchQuery() == "" || utils.Match(torrentItem.Filename, config.SearchQuery()) {
			_, err := fmt.Fprintln(v, torrentItem.Filename)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

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

func GetSelectedTorrent(v *gocui.View) (torrent models.Torrent, cursorPosition int, err error) {
	_, cy := v.Cursor()
	var emptyTorrent models.Torrent
	if cy < 0 {
		return emptyTorrent, cy, fmt.Errorf("cursor is off-screen or uninitialized")
	}
	if cy >= len(data.UserDownloads) {
		return emptyTorrent, cy, fmt.Errorf("cursor index %d out of bounds (max %d)", cy, len(data.UserDownloads)-1)
	}
	return data.UserDownloads[cy], cy, nil
}

// TODO
// this wont work if the viewport is too small to show the entire id
func GetSelectedActiveDownload(v *gocui.View) (string, error) {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	return line, err
}

func GetSelectedTorrentFile(v *gocui.View) (models.TorrentFileDetailed, error) {
	_, cy := v.Cursor()
	line, err := v.Line(cy)

	if err != nil {
		return models.TorrentFileDetailed{}, fmt.Errorf("unable to get selected line: %w", err)
	}

	torrentFile, ok := data.FilesMap[line]

	if !ok {
		return models.TorrentFileDetailed{}, fmt.Errorf("no download item found for selected line")
	}
	return torrentFile, nil
}

func PopulateViews(g *gocui.Gui) {
	torrentsView := GetView(g, ViewTorrents)
	torrentsView.Clear()

	// temporary. needs a centralized way that uses the error model
	if len(data.UserDownloads) == 0 {
		if len(data.UserDownloads) == 0 && !warnedNoTorrents {
			UpdateUILog(g, "API returned no torrents, is your API token correct?", nil)
			warnedNoTorrents = true
		}
	}

	activeView := GetView(g, ViewActiveTorrents)
	activeView.Clear()

	// Add active downloads from the API
	for index := range len(data.UserDownloads) {
		item := data.UserDownloads[index]

		if item.Status == "queued" || item.Status == "downloading" {
			fmt.Fprintln(activeView, item.ID)
		} else if item.Status == "downloaded" {
			fmt.Fprintln(torrentsView, strings.TrimSpace(item.Filename))

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

	torrentItem, ok := data.UserDownloads[cy]
	if !ok {
		_, err = fmt.Fprint(mainView, "No details found.")
		if err != nil {
			return err
		}
		return nil
	}

	_, err = fmt.Fprint(mainView, GenerateDetailsString(torrentItem))
	if err != nil {
		return err
	}
	return nil
}

func CopyDownloadLink(g *gocui.Gui, v *gocui.View) error {
	item, err := GetSelectedTorrentFile(v)
	if err != nil {
		return err
	}

	if err := clipboard.WriteAll(item.Download); err != nil {
		UpdateUILog(g, "Failed to copy download link:", err)
		return err
	}

	UpdateUILog(g, fmt.Sprintf("Copied download link for %s", item.Filename), nil)
	return nil
}

func SearchKeyPress(g *gocui.Gui, v *gocui.View) error {
	config.SetSearchQuery(v.Buffer())

	if err := RenderList(g); err != nil {
		return err
	}

	torrentsView, _ := g.View(ViewTorrents)
	if err := UpdateDetails(g, torrentsView); err != nil {
		return err
	}

	_, err := g.SetCurrentView(ViewTorrents)
	return err
}

func ShowTorrentFiles(g *gocui.Gui, v *gocui.View, fileMap map[string]models.TorrentFileDetailed) {

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

		_, err := g.SetCurrentView(ViewDetails)
		if err != nil {
			logs.LogEvent(err)
		}
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
