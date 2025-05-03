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
	"time"

	"lazydebrid/internal/config"
	"lazydebrid/internal/models"

	"github.com/jroimartin/gocui"
)

const (
	downloadsURL   = "https://api.real-debrid.com/rest/1.0/downloads?page=1&limit=4990"
	addMagnetURL   = "https://api.real-debrid.com/rest/1.0/torrents/addMagnet"
	statusURL      = "https://api.real-debrid.com/rest/1.0/torrents/info/"
	deleteURL      = "https://api.real-debrid.com/rest/1.0/torrents/delete/"
	selectFilesURL = "https://api.real-debrid.com/rest/1.0/torrents/selectFiles/"
)

var (
	UserDownloads   []models.DebridDownload
	DownloadMap     = make(map[string]models.DebridDownload)
	ActiveDownloads []models.ActiveDownload
)

func DownloadFile(torrentItem models.DebridDownload) bool {
	out, err := os.Create(fmt.Sprintf("%s%s", config.UserDownloadPath, torrentItem.Filename))
	defer out.Close()

	if err != nil {
		log.Println(err)
	}
	resp, err := http.Get(torrentItem.Download)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	n, err := io.Copy(out, resp.Body)
	if err == nil {
		log.Println(n)
		return true
	}
	return false
}

func GetTorrentStatus(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return err
	}

	statusURL := fmt.Sprintf("https://api.real-debrid.com/rest/1.0/torrents/info/%s", line)
	req, err := http.NewRequest("GET", statusURL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.UserApiToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	infoView, _ := g.View("info")
	now := time.Now().Format("02 Jan 2006 15:04:00")

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(infoView, "[%s] Failed to fetch status: %s\n", now, msg)
		return nil
	}

	var info models.TorrentStatus
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		fmt.Fprintf(infoView, "[%s] Failed to parse status\n", now)
		return err
	}

	fmt.Fprintf(infoView,
		"[%s]\nStatus for %s:\n  Status: %s\n  Progress: %d%%\n  Added: %s\n  Files: %d\n\n",
		now, info.Filename, info.Status, info.Progress, info.Added, len(info.Files),
	)

	return nil
}
func DeleteTorrent(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return err
	}
	deleteURL := fmt.Sprintf("https://api.real-debrid.com/rest/1.0/torrents/delete/%s", line)
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.UserApiToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	infoView, _ := g.View("info")
	now := time.Now().Format("02 Jan 2006 15:04:00")
	if resp.StatusCode == http.StatusNoContent {
		fmt.Fprintf(infoView, "[%s] Deleted torrent: %s", now, line)
	} else {
		msg, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(infoView, "[%s] Failed to delete %s: %s", now, line, msg)
	}
	return nil
}
func AddFilesToDebrid(downloadID string) bool {
	data := url.Values{}
	data.Set("files", "all")

	selectFilesURL := fmt.Sprintf("https://api.real-debrid.com/rest/1.0/torrents/selectFiles/%s", downloadID)
	req, err := http.NewRequest("POST", selectFilesURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.UserApiToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	log.Printf("Status code: %d", resp.StatusCode)
	if resp.StatusCode != 204 && resp.StatusCode != 202 {
		msg, _ := io.ReadAll(resp.Body)
		log.Println(fmt.Errorf("HTTP %d: %s", resp.StatusCode, msg))
		return false
	}

	return true
}

func SendLinkToAPI(magnetLink string) (string, error) {
	data := url.Values{}
	data.Set("magnet", magnetLink)

	req, err := http.NewRequest("POST", addMagnetURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.UserApiToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		msg, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, msg)
	}

	var result models.ActiveDownload
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding addMagnet response: %w", err)
	}

	log.Printf("Magnet added. Torrent ID: %s", result.ID)
	ActiveDownloads = append(ActiveDownloads, result)
	success := AddFilesToDebrid(result.ID)

	if success {
		return result.ID, nil
	}
	return "", fmt.Errorf("Could not add files")
}

func GetUserTorrents() map[string]models.DebridDownload {

	client := &http.Client{}
	req, err := http.NewRequest("GET", downloadsURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.UserApiToken))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &UserDownloads)
	if err != nil {
		log.Fatal(err)
	}

	for _, torrentItem := range UserDownloads {
		DownloadMap[torrentItem.Filename] = torrentItem
	}

	return DownloadMap
}
