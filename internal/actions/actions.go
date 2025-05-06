package actions

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"lazydebrid/internal/api"
	"lazydebrid/internal/config"
	"lazydebrid/internal/data"
	"lazydebrid/internal/logs"
	"lazydebrid/internal/models"
	"lazydebrid/internal/views"

	"github.com/jroimartin/gocui"
)

func DeleteTorrent(g *gocui.Gui, v *gocui.View) error {
	torrentID, err := views.GetSelectedLine(v)
	if err != nil || strings.TrimSpace(torrentID) == "" {
		return fmt.Errorf("no torrent selected")
	}

	req, err := api.NewRequest("DELETE", api.TorrentsDeleteURL+torrentID, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	_, err = api.DoRequest(req)
	if err != nil {
		views.UpdateUILog(g, fmt.Sprintf("Failed to delete torrent: %s\nError: %s", torrentID, err), false, nil)
		return err
	}

	// if body is empty, it succeeded
	data.ActiveDownloads = RemoveItem(data.ActiveDownloads, torrentID)
	views.UpdateUILog(g, fmt.Sprintf("Deleted torrent: %s", torrentID), true, nil)

	return nil
}
func AddFilesToDebrid(downloadID string) bool {
	form := url.Values{"files": {"all"}}
	_, err := api.PostForm(api.TorrentsSelectFilesURL+downloadID, form)
	if err != nil {
		logs.LogEvent(fmt.Errorf("select files API call failed: %w", err))
		return false
	}
	return true
}

func SendLinkToAPI(magnetLink string) (string, error) {
	payload := url.Values{"magnet": {magnetLink}}

	body, err := api.PostForm(api.TorrentsAddMagnetURL, payload)
	if err != nil {
		return "", fmt.Errorf("failed to add magnet: %w", err)
	}

	var download models.ActiveDownload
	if err := json.Unmarshal(body, &download); err != nil {
		return "", fmt.Errorf("failed to parse add magnet response: %w", err)
	}

	data.ActiveDownloads = append(data.ActiveDownloads, download)
	log.Printf("ActiveDownloads now: %d entries", len(data.ActiveDownloads))

	if ok := AddFilesToDebrid(download.ID); !ok {
		return "", fmt.Errorf("magnet added but file selection failed")
	}

	return download.ID, nil
}

func GetTorrentStatus(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil || strings.TrimSpace(line) == "" {
		return fmt.Errorf("no torrent selected")
	}
	torrentID := strings.TrimSpace(line)

	req, err := api.NewRequest("GET", api.TorrentsStatusURL+torrentID, nil)
	if err != nil {
		return fmt.Errorf("failed to create status request: %w", err)
	}

	body, err := api.DoRequest(req)
	if err != nil {
		views.UpdateUILog(g, fmt.Sprintf("Failed to get torrent status: %v", err), false, nil)
		return err
	}

	var status models.Torrent
	if err := json.Unmarshal(body, &status); err != nil {
		logs.LogEvent(err)
		views.UpdateUILog(g, "Failed to decode torrent status", false, err)
		return fmt.Errorf("error decoding status: %w", err)
	}

	views.UpdateUILog(g, fmt.Sprintf(
		"\nStatus for %s:\n  Status: %s\n  Progress: %d%%\n  Added: %s\n  Files: %d\n\n",
		status.Filename, status.Status, status.Progress, status.Added, len(status.Files)), true, nil)

	return nil
}

func DownloadFile(torrent models.Download) bool {
	path := filepath.Join(config.DownloadPath(), torrent.Filename)

	resp, err := http.Get(torrent.Download)
	if err != nil {
		logs.LogEvent(fmt.Errorf("failed to GET %s: %w", torrent.Download, err))
		return false
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		logs.LogEvent(fmt.Errorf("failed to create file %s: %w", path, err))
		return false
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		logs.LogEvent(fmt.Errorf("failed to write to file %s: %w", path, err))
		return false
	}

	logs.LogEvent(fmt.Errorf("downloaded: %s", path))
	return true
}

func GetTorrentContents(g *gocui.Gui, v *gocui.View) map[string]models.Download {
	id, err := views.GetSelectedTorrentID(v)
	if err != nil {
		logs.LogEvent(fmt.Errorf("selection error: %w", err))
		views.UpdateUILog(g, "No torrent selected", false, nil)
		return nil
	}

	torrent, ok := data.DownloadMap[id]
	if !ok {
		msg := fmt.Sprintf("No torrent found for ID: %s", id)
		logs.LogEvent(fmt.Errorf(msg))
		views.UpdateUILog(g, msg, false, nil)
		return nil
	}

	files := make(map[string]models.Download)
	var errors []string

	for _, link := range torrent.Links {
		file, err := api.UnrestrictLink(link)
		if err != nil {
			logs.LogEvent(err)
			errors = append(errors, err.Error())
			continue
		}
		files[file.Filename] = file
	}

	data.FilesMap = files

	if len(errors) > 0 {
		views.UpdateUILog(g, strings.Join(errors, "; "), false, nil)
	}

	return files
}

func GetUserTorrents() map[string]models.Torrent {
	result := make(map[string]models.Torrent)

	req, err := api.NewRequest("GET", api.TorrentsURL, nil)
	if err != nil {
		logs.LogEvent(fmt.Errorf("failed to create request for user torrents: %w", err))
		return result
	}

	body, err := api.DoRequest(req)
	if err != nil {
		logs.LogEvent(fmt.Errorf("failed to fetch user torrents: %w", err))
		return result
	}

	var list []models.Torrent
	if err := json.Unmarshal(body, &list); err != nil {
		logs.LogEvent(fmt.Errorf("failed to parse torrent list: %w", err))
		return result
	}

	data.TorrentLineIndex = data.TorrentLineIndex[:0] // clear before repopulating
	for _, item := range list {
		data.TorrentLineIndex = append(data.TorrentLineIndex, item.Filename)
		result[item.Filename] = item
	}

	data.DownloadMap = result
	data.UserDownloads = list

	return result
}
