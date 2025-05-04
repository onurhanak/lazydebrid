package actions

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"lazydebrid/internal/config"
	"lazydebrid/internal/logs"
	"lazydebrid/internal/logui"
	"lazydebrid/internal/models"
	"lazydebrid/internal/views"

	"github.com/jroimartin/gocui"
)

const (
	baseURL                = "https://api.real-debrid.com/rest/1.0"
	torrentsEndpointURL    = baseURL + "/torrents"
	downloadsURL           = baseURL + "/downloads?page=1&limit=4990"
	torrentsURL            = baseURL + "/torrents/"
	torrentsAddMagnetURL   = torrentsEndpointURL + "/addMagnet/"
	torrentsStatusURL      = torrentsEndpointURL + "/info/"
	torrentsDeleteURL      = torrentsEndpointURL + "/delete/"
	torrentsSelectFilesURL = torrentsEndpointURL + "/selectFiles/"
)

var (
	UserDownloads   []models.Torrent
	DownloadMap     = make(map[string]models.Torrent)
	ActiveDownloads []models.ActiveDownload
	FilesMap        = make(map[string]models.Download)
)

type LineMapping struct {
	ID string
}

var TorrentLineIndex []string

func newRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		logs.LogEvent(err)
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.APIToken()))
	return req, nil
}

func readResponse(resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, fmt.Errorf("response is nil")
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func doRequest(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.LogEvent(err)
		return nil, err
	}
	return resp, nil
}

func DeleteTorrent(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil || line == "" {
		return fmt.Errorf("no torrent selected")
	}

	req, err := newRequest("DELETE", torrentsDeleteURL+line, nil)
	if err != nil {
		return err
	}

	resp, err := doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	infoView := views.GetView(g, views.ViewInfo)
	now := logs.GetNow()

	if resp.StatusCode == http.StatusNoContent {
		logui.LogInfo(infoView, now, fmt.Sprintf("Deleted torrent: %s", line))
	} else {
		msg, _ := io.ReadAll(resp.Body)
		logui.LogError(infoView, now, fmt.Sprintf("Failed to delete %s: %s", line, msg), nil)
	}

	return nil
}

func AddFilesToDebrid(downloadID string) bool {
	data := url.Values{"files": {"all"}}
	req, err := newRequest("POST", torrentsSelectFilesURL+downloadID, strings.NewReader(data.Encode()))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := doRequest(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusAccepted {
		msg, _ := io.ReadAll(resp.Body)

		logs.LogEvent(fmt.Errorf("failed to select files: HTTP %d: %s", resp.StatusCode, msg))
		return false
	}

	return true
}

func SendLinkToAPI(magnetLink string) (string, error) {
	data := url.Values{"magnet": {magnetLink}}
	req, err := newRequest("POST", torrentsAddMagnetURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := doRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		msg, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to add magnet: HTTP %d: %s", resp.StatusCode, msg)
	}

	var result models.ActiveDownload
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	ActiveDownloads = append(ActiveDownloads, result)
	log.Printf("ActiveDownloads now: %d entries", len(ActiveDownloads))
	if !AddFilesToDebrid(result.ID) {
		return "", fmt.Errorf("magnet added but failed to select files")
	}
	return result.ID, nil
}

func GetTorrentStatus(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil || strings.TrimSpace(line) == "" {
		return fmt.Errorf("no torrent selected")
	}
	torrentID := strings.TrimSpace(line)

	req, err := newRequest("GET", torrentsStatusURL+torrentID, nil)
	if err != nil {
		return err
	}

	resp, err := doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logui.LogError(views.GetView(g, views.ViewInfo), logs.GetNow(),
			fmt.Sprintf("Failed to get torrent status: HTTP %d\n%s", resp.StatusCode, string(body)), nil)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var status models.Torrent
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		logs.LogEvent(err)
		logui.LogError(views.GetView(g, views.ViewInfo), logs.GetNow(), "Failed to decode torrent status", err)
		return err
	}

	infoView := views.GetView(g, views.ViewInfo)
	now := logs.GetNow()
	_, err = fmt.Fprintf(infoView,
		"[%s]\nStatus for %s:\n  Status: %s\n  Progress: %d%%\n  Added: %s\n  Files: %d\n\n",
		now, status.Filename, status.Status, status.Progress, status.Added, len(status.Files),
	)
	return err
}
func DownloadFile(torrent models.Download) bool {
	path := fmt.Sprintf("%s%s", config.DownloadPath(), torrent.Filename)
	out, err := os.Create(path)
	if err != nil {
		logs.LogEvent(err)
		return false
	}
	defer out.Close()

	resp, err := http.Get(torrent.Download)
	if err != nil {
		logs.LogEvent(err)
		return false
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logs.LogEvent(err)
		return false
	}

	logs.LogEvent(fmt.Errorf("downloaded: %s", path))
	return true
}

func GetSelectedTorrentID(v *gocui.View) (string, error) {
	_, cy := v.Cursor()
	if cy < 0 || cy >= len(TorrentLineIndex) {
		return "", fmt.Errorf("invalid cursor line index: %d", cy)
	}
	return TorrentLineIndex[cy], nil
}

func GetTorrentContents(g *gocui.Gui, v *gocui.View) map[string]models.Download {
	id, err := GetSelectedTorrentID(v)
	if err != nil {
		log.Println("Selection error:", err)
		return nil
	}

	torrent, ok := DownloadMap[id]
	if !ok {
		log.Printf("No torrent found for ID: %s", id)
		return nil
	}

	var torrentFile models.Download
	var errorMessage struct {
		Error     string `json:"error"`
		ErrorCode int    `json:"error_code"`
	}

	var errorLog []string
	FilesMap = make(map[string]models.Download)

	for _, link := range torrent.Links {
		data := url.Values{"link": {link}}
		req, _ := newRequest("POST", baseURL+"/unrestrict/link/", strings.NewReader(data.Encode()))
		resp, _ := doRequest(req)
		response, _ := readResponse(resp)

		// does not work
		if err := json.Unmarshal(response, &torrentFile); err != nil {
			logs.LogEvent(err)
			if err = json.Unmarshal(response, &errorMessage); err == nil {
				logs.LogEvent(fmt.Errorf("API error: %s (code %d)", errorMessage.Error, errorMessage.ErrorCode))
				errorLog = append(errorLog, errorMessage.Error)
			} else {
				errorLog = append(errorLog, "Unmarshal failed and no API error given")
			}
			log.Println("log")
			log.Println(errorLog)
			continue
		}

		if torrentFile.Filename != "" {
			FilesMap[torrentFile.Filename] = torrentFile
		}
	}

	if len(errorLog) > 0 {
		msg := strings.Join(errorLog, "; ")
		log.Println(msg)
		g.Update(func(g *gocui.Gui) error {
			logui.UpdateUILog(g, msg)
			return nil
		})
	}

	return FilesMap
}

func GetUserTorrents() map[string]models.Torrent {
	result := make(map[string]models.Torrent)

	req, err := newRequest("GET", torrentsURL, nil)
	if err != nil {
		logs.LogEvent(err)
		return result
	}

	resp, err := doRequest(req)
	if err != nil {
		logs.LogEvent(err)
		return result
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logs.LogEvent(err)
		return result
	}

	var list []models.Torrent
	if err := json.Unmarshal(body, &list); err != nil {
		logs.LogEvent(err)
		return result
	}

	for _, item := range list {
		TorrentLineIndex = append(TorrentLineIndex, item.Filename)
		result[item.Filename] = item
	}
	DownloadMap = result
	UserDownloads = list

	return result
}

func DumpTorrentsToJSON(torrents map[string]models.Torrent, filename string) error {
	data, err := json.MarshalIndent(torrents, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal torrents: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}
