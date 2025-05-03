package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/jroimartin/gocui"
)

// LOG
func init() {
	logFile, err := os.OpenFile("lazydebrid.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Could not open log file:", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

const (
	downloadsURL   = "https://api.real-debrid.com/rest/1.0/downloads?page=1&limit=5000"
	addMagnetURL   = "https://api.real-debrid.com/rest/1.0/torrents/addMagnet"
	statusURL      = "https://api.real-debrid.com/rest/1.0/torrents/info/"
	deleteURL      = "https://api.real-debrid.com/rest/1.0/torrents/delete/"
	selectFilesURL = "https://api.real-debrid.com/rest/1.0/torrents/selectFiles/"
	configFile     = "lazyDebrid.json"
)

var (
	userDownloads   []DebridDownload
	userSettings    = make(map[string]string)
	downloadMap     = make(map[string]DebridDownload)
	activeDownloads []ActiveDownload
	views           = []string{"search", "torrents", "details", "activeTorrents"}
	currentViewIdx  int

	userApiToken     string
	userDownloadPath string
	searchQuery      string

	userConfigPath, _   = os.UserConfigDir()
	lazyDebridConfig    = filepath.Join(userConfigPath, "lazyDebrid.json")
	mainKeys            = "Switch Pane: <Tab> | Download Path: <^P> | API Key: <^X> | Quit: <^Q>"
	torrentsKeys        = "Up: <k> | Down: <j> | Add Magnet: <a> | Copy link: <y> | Download: <d> | Keybindings: <?>"
	activeDownloadsKeys = "Up: <k> | Down: <j> | Status: <s> | Delete <d> | Keybindings: <?>"
	searchKeys          = "Search: <Enter> | Keybindings: <?>"
)

// MODELS
type DebridDownload struct {
	Id         string `json:"id"`
	Filename   string `json:"filename"`
	MimeType   string `json:"mimeType"`
	Filesize   int64  `json:"filesize"`
	Link       string `json:"link"`
	Host       string `json:"host"`
	HostIcon   string `json:"host_icon"`
	Chunks     int64  `json:"chunks"`
	Download   string `json:"download"`
	Streamable int64  `json:"streamable"`
	Generated  string `json:"generated"`
}

type TorrentStatus struct {
	ID               string        `json:"id"`
	Filename         string        `json:"filename"`
	OriginalFilename string        `json:"original_filename"`
	Hash             string        `json:"hash"`
	Bytes            int64         `json:"bytes"`
	OriginalBytes    int64         `json:"original_bytes"`
	Host             string        `json:"host"`
	Split            int           `json:"split"`
	Progress         int           `json:"progress"`
	Status           string        `json:"status"`
	Added            string        `json:"added"`
	Files            []TorrentFile `json:"files"`
	Links            []string      `json:"links"`
}

type TorrentFile struct {
	ID       int    `json:"id"`
	Path     string `json:"path"`
	Bytes    int64  `json:"bytes"`
	Selected int    `json:"selected"`
}

type ActiveDownload struct {
	ID  string `json:"id"`
	URI string `json:"uri"`
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	currentViewIdx = (currentViewIdx + 1) % len(views)
	name := views[currentViewIdx]
	_, err := g.SetCurrentView(name)

	keysView, _ := g.View("footer")
	keysView.Clear()
	if name == "torrents" {
		fmt.Fprint(keysView, torrentsKeys)
	} else if name == "search" {
		fmt.Fprint(keysView, searchKeys)
	} else if name == "activeTorrents" {
		fmt.Fprint(keysView, activeDownloadsKeys)
	} else if name == "details" {
		fmt.Fprint(keysView, "")
	}

	return err
}

// CONFIG

func configPath() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, configFile)
}

func loadUserSettings() error {
	data, err := os.ReadFile(configPath())
	if err == nil {
		_ = json.Unmarshal(data, &userSettings)
	}
	userApiToken = strings.TrimSpace(userSettings["apiToken"])
	userDownloadPath = strings.TrimSpace(userSettings["downloadPath"])
	return nil
}

func saveSetting(key, value string) error {
	value = strings.TrimSpace(value)
	data, _ := os.ReadFile(configPath())
	_ = json.Unmarshal(data, &userSettings)
	userSettings[key] = value
	content, err := json.MarshalIndent(userSettings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), content, 0644)
}

func showModal(g *gocui.Gui, name, title string, content string, onSubmit func(string)) error {
	maxX, maxY := g.Size()
	w, h := maxX/2, 5
	x0 := (maxX - w) / 2
	y0 := (maxY - h) / 2
	x1 := x0 + w
	y1 := y0 + h

	if v, err := g.SetView(name, x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = title
		if name != "help" {
			v.Title = title
			v.Editable = true
			v.Wrap = true
		} else {
			v.Editable = false
			v.Wrap = false
			fmt.Fprint(v, content)
		}
		g.DeleteKeybindings(name)
		_, _ = g.SetCurrentView(name)
		g.SetKeybinding(name, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			input := strings.TrimSpace(v.Buffer())
			g.DeleteView(name)
			g.SetCurrentView("torrents")
			onSubmit(input)
			return nil
		})
	}
	return nil
}

func showSetPathModal(g *gocui.Gui, v *gocui.View) error {
	return showModal(g, "setPath", "Set Download Path", "", func(input string) {
		_ = saveSetting("downloadPath", input)
	})
}

func showSetTokenModal(g *gocui.Gui, v *gocui.View) error {
	return showModal(g, "setToken", "Set API Token", "", func(input string) {
		_ = saveSetting("apiToken", input)
	})
}

func showAddMagnetModal(g *gocui.Gui, v *gocui.View) error {
	return showModal(g, "addMagnet", "Add Magnet Link", "", func(input string) {
		_, _ = sendLinkToAPI(input)
	})
}

func showHelpModal(g *gocui.Gui, v *gocui.View) error {
	content := "TAB: Switch | ↑↓: Navigate | ENTER: Download | /: Search\n^A: Add Magnet | ^C: Copy Link | ^P: Set Path\n^X: Set API Key | ^Q: Quit"

	return showModal(g, "help", "Shortcuts", content, func(string) {})
}

func saveApiToken(apiToken string) bool {
	if content, err := os.ReadFile(lazyDebridConfig); err == nil {
		_ = json.Unmarshal(content, &userSettings)
	} else if !errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}
	apiToken = strings.ReplaceAll(apiToken, "\n", "")
	userSettings["apiToken"] = apiToken

	data, err := json.MarshalIndent(userSettings, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(lazyDebridConfig, data, 0644); err != nil {
		log.Fatal(err)
	}

	return true
}

func saveDownloadPath(downloadPath string) bool {

	downloadPath = strings.ReplaceAll(downloadPath, "\n", "")
	if content, err := os.ReadFile(lazyDebridConfig); err == nil {
		_ = json.Unmarshal(content, &userSettings)
	} else if !errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	userSettings["downloadPath"] = downloadPath

	data, err := json.MarshalIndent(userSettings, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(lazyDebridConfig, data, 0644); err != nil {
		log.Fatal(err)
	}

	return true
}

func setDownloadPath(g *gocui.Gui, v *gocui.View) error {
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	maxX, maxY := g.Size()
	if v, err := g.SetView("setDownloadPath", maxX/4, maxY/4, maxX*3/4, maxY*3/10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Set download path"
		v.Wrap = true
		v.Editable = true
		v.Clear()
		if _, err := g.SetCurrentView("setDownloadPath"); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("setDownloadPath", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		downloadPath := v.Buffer()
		saveDownloadPath(downloadPath)
		g.DeleteView("setDownloadPath")
		g.SetCurrentView("torrents")
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func setApiToken(g *gocui.Gui, v *gocui.View) error {
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	maxX, maxY := g.Size()
	if v, err := g.SetView("setApiToken", maxX/4, maxY/4, maxX*3/4, maxY*3/10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Set Api Token"
		v.Wrap = true
		v.Editable = true
		v.Clear()
		if _, err := g.SetCurrentView("setApiToken"); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("setApiToken", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		downloadPath := v.Buffer()
		saveApiToken(downloadPath)
		g.DeleteView("setApiToken")
		g.SetCurrentView("torrents")
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// TORRENTS

func getUserTorrents() map[string]DebridDownload {

	client := &http.Client{}
	req, err := http.NewRequest("GET", downloadsURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userApiToken))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &userDownloads)
	if err != nil {
		log.Fatal(err)
	}

	for _, torrentItem := range userDownloads {
		downloadMap[torrentItem.Filename] = torrentItem
	}

	log.Println(downloadMap)
	return downloadMap
}

// UTILS
func match(filename, query string) bool {
	return strings.Contains(strings.ToLower(filename), strings.ToLower(query))
}

func renderList(g *gocui.Gui) error {
	v, err := g.View("torrents")
	if err != nil {
		return err
	}
	v.Clear()
	v.SetCursor(0, 0)

	for _, torrentItem := range userDownloads {
		if searchQuery == "" || match(torrentItem.Filename, searchQuery) {
			fmt.Fprintln(v, torrentItem.Filename)
		}
	}

	updateDetails(g, v)

	return nil
}

func getCurrentLine(v *gocui.View) (string, error) {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	return line, err
}

func generateDetailsString(torrentItem DebridDownload) string {
	detailsString := fmt.Sprintf(
		"ID: %s\nFilename: %s\nMIME Type: %s\nFilesize: %d bytes\nLink: %s\nDownload: %s\nStreamable: %d",
		torrentItem.Id,
		torrentItem.Filename,
		torrentItem.MimeType,
		torrentItem.Filesize,
		torrentItem.Link,
		torrentItem.Download,
		torrentItem.Streamable,
	)
	return detailsString
}

// HANDLERS

func searchKeyPress(g *gocui.Gui, v *gocui.View) error {
	searchQuery = strings.TrimSpace(v.Buffer())
	renderList(g)
	g.SetCurrentView("torrents")
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy+1); err != nil {
		ox, oy := v.Origin()
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}

	return updateDetails(g, v)
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	if cy > 0 {
		if err := v.SetCursor(cx, cy-1); err != nil {
			return err
		}
	}
	return updateDetails(g, v)
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// ACTIONS

func getTorrentStatus(g *gocui.Gui, v *gocui.View) error {
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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userApiToken))

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

	var info TorrentStatus
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

func deleteTorrent(g *gocui.Gui, v *gocui.View) error {
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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userApiToken))

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

func downloadSelected(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return err
	}
	downloadItem := downloadMap[line]
	v, _ = g.View("info")
	now := time.Now().Format("02 Jan 2006 15:04:00")

	fmt.Fprint(v, fmt.Sprintf("[%s] Downloading %s to %s", now, downloadItem.Filename, userDownloadPath))
	go func(torrentItem DebridDownload) {
		if downloadFile(torrentItem) {
			fmt.Fprint(v, fmt.Sprintf("Downloaded %s to %s", torrentItem.Filename, userDownloadPath))
		}
	}(downloadItem)
	return nil
}

func copyDownloadLink(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return err
	}
	selectedItem := downloadMap[line]
	err = clipboard.WriteAll(selectedItem.Download)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func downloadFile(torrentItem DebridDownload) bool {
	out, err := os.Create(fmt.Sprintf("%s%s", userDownloadPath, torrentItem.Filename))
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

func focusSearchBar(g *gocui.Gui, v *gocui.View) error {
	g.SetCurrentView("search")
	return nil
}

func updateDetails(g *gocui.Gui, v *gocui.View) error {
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

	torrentItem, ok := downloadMap[strings.TrimSpace(line)]
	if !ok {
		fmt.Fprint(mainView, "No details found.")
		return nil
	}

	fmt.Fprint(mainView, generateDetailsString(torrentItem))
	return nil
}

func addFilesToDebrid(downloadID string) bool {
	data := url.Values{}
	data.Set("files", "all")

	selectFilesURL := fmt.Sprintf("https://api.real-debrid.com/rest/1.0/torrents/selectFiles/%s", downloadID)
	req, err := http.NewRequest("POST", selectFilesURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userApiToken))

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

func sendLinkToAPI(magnetLink string) (string, error) {
	data := url.Values{}
	data.Set("magnet", magnetLink)

	req, err := http.NewRequest("POST", addMagnetURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userApiToken))

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

	var result ActiveDownload
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding addMagnet response: %w", err)
	}

	log.Printf("Magnet added. Torrent ID: %s", result.ID)
	activeDownloads = append(activeDownloads, result)
	success := addFilesToDebrid(result.ID)

	if success {
		return result.ID, nil
	}
	return "", fmt.Errorf("Could not add files")
}

func addMagnetLink(g *gocui.Gui, v *gocui.View) error {
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	maxX, maxY := g.Size()
	_, err := g.SetView("addMagnet", maxX/4, maxY/4, maxX*3/4, maxY*3/4)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	v, _ = g.View("addMagnet")
	v.Title = "Add Magnet Link"
	v.Wrap = true
	v.Editable = true
	v.Clear()
	g.SetCurrentView("addMagnet")

	if err := g.SetKeybinding("addMagnet", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		magnetLink := strings.TrimSpace(v.Buffer())
		infoView, _ := g.View("info")
		now := time.Now().Format("02 Jan 2006 15:04:00")

		if magnetLink == "" {
			fmt.Fprintf(infoView, "[%s] Error: Empty magnet link\n", now)
			return nil
		}

		downloadID, err := sendLinkToAPI(magnetLink)
		if err != nil {
			fmt.Fprintf(infoView, "[%s] Failed to add magnet: %v\n", now, err)
			return nil
		}
		fmt.Fprintf(infoView, "[%s] Magnet added: %s\n", now, downloadID)

		success := addFilesToDebrid(downloadID)

		if success {
			fmt.Fprintf(infoView, "[%s] All files selected for download: %s\n", now, downloadID)
		} else {
			fmt.Fprintf(infoView, "[%s] Failed to select files for %s\n", now, downloadID)
		}

		g.DeleteView("addMagnet")
		g.SetCurrentView("torrents")
		return nil
	}); err != nil {
		if !strings.Contains(err.Error(), "duplicate") {
			return err
		}
	}

	return nil
}

func deleteCurrentView(g *gocui.Gui, v *gocui.View) error {
	currentView := g.CurrentView()
	log.Println(currentView.Name())
	if currentView == nil {
		return nil
	}

	if currentView.Name() != "torrents" && currentView.Name() != "details" && currentView.Name() != "info" && currentView.Name() != "footer" && currentView.Name() != "activeTorrents" && currentView.Name() != "search" {
		g.DeleteView(currentView.Name())
		g.SetCurrentView("torrents")
	}

	return nil
}

func showControls(g *gocui.Gui, v *gocui.View) error {
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
		fmt.Fprint(v, controlsString)
		g.SetCurrentView("controls")
	}

	return nil
}

// BINDINGS

func keybindings(g *gocui.Gui) error {
	bind := func(viewname string, key interface{}, mod gocui.Modifier, handler func(*gocui.Gui, *gocui.View) error) {
		if err := g.SetKeybinding(viewname, key, mod, handler); err != nil {
			log.Fatalf("binding failed: %v", err)
		}
	}

	bind("activeTorrents", 'd', gocui.ModNone, deleteTorrent)
	bind("activeTorrents", 's', gocui.ModNone, getTorrentStatus)
	bind("activeTorrents", 'j', gocui.ModNone, cursorDown)
	bind("activeTorrents", 'k', gocui.ModNone, cursorUp)
	bind("torrents", 'j', gocui.ModNone, cursorDown)
	bind("torrents", 'k', gocui.ModNone, cursorUp)
	bind("", gocui.KeyArrowDown, gocui.ModNone, cursorDown)
	bind("", gocui.KeyArrowUp, gocui.ModNone, cursorUp)
	bind("torrents", gocui.KeyEnter, gocui.ModNone, downloadSelected)
	bind("torrents", '/', gocui.ModNone, focusSearchBar)
	bind("details", '/', gocui.ModNone, focusSearchBar)
	bind("search", gocui.KeyEnter, gocui.ModNone, searchKeyPress)
	bind("", gocui.KeyCtrlC, gocui.ModNone, copyDownloadLink)
	bind("", gocui.KeyCtrlD, gocui.ModNone, downloadSelected)
	bind("", gocui.KeyCtrlA, gocui.ModNone, addMagnetLink)
	bind("", gocui.KeyCtrlP, gocui.ModNone, showSetPathModal)
	bind("", gocui.KeyCtrlX, gocui.ModNone, showSetTokenModal)
	bind("", gocui.KeyCtrlQ, gocui.ModNone, quit)
	bind("", gocui.KeyTab, gocui.ModNone, nextView)
	bind("torrents", '?', gocui.ModNone, showHelpModal)
	return nil
}

// LAYOUT

func layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()
	splitX := (maxX * 4) / 10
	infoHeight := (maxY - 3) / 4

	detailsTop := 3
	detailsBottom := detailsTop + infoHeight
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen

	activeTop := detailsBottom + 1
	activeBottom := activeTop + infoHeight

	infoTop := activeBottom + 1
	infoBottom := maxY - 4
	if v, err := g.SetView("search", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Search"
		v.Editable = true
	}

	if torrentsView, err := g.SetView("torrents", 0, 3, splitX, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		torrentsView.Title = "Downloads"
		torrentsView.Highlight = true
		torrentsView.Wrap = false
		torrentsView.SelFgColor = gocui.ColorGreen

		for _, item := range userDownloads {
			fmt.Fprintln(torrentsView, item.Filename)
		}
		torrentsView.SetCursor(0, 0)
		g.SetCurrentView("torrents")
		updateDetails(g, torrentsView)
	}

	if mainView, err := g.SetView("details", splitX+1, detailsTop, maxX-1, detailsBottom); err != nil && err != gocui.ErrUnknownView {
		return err
	} else if err == nil {
		mainView.Title = "Torrent Details"
		mainView.Wrap = true

		if torrentsView, err := g.View("torrents"); err == nil {
			updateDetails(g, torrentsView)
		}
	}

	if activeTorrentsView, err := g.SetView("activeTorrents", splitX+1, activeTop, maxX-1, activeBottom); err != nil && err != gocui.ErrUnknownView {
		return err
	} else if err == nil {
		activeTorrentsView.Title = "Active Downloads"
		activeTorrentsView.Highlight = true
		activeTorrentsView.Wrap = false
		activeTorrentsView.SelFgColor = gocui.ColorGreen
		activeTorrentsView.Clear()
		for _, item := range activeDownloads {
			fmt.Fprintln(activeTorrentsView, item.ID)
		}
		activeTorrentsView.SetCursor(0, 0)
	}

	if infoView, err := g.SetView("info", splitX+1, infoTop, maxX-1, infoBottom); err != nil && err != gocui.ErrUnknownView {
		return err
	} else if err == nil {
		infoView.Title = "Log"
		infoView.Wrap = true
		infoView.Autoscroll = true
	}

	if footerView, err := g.SetView("footer", 0, infoBottom+1, maxX-1, infoBottom+3); err != nil && err != gocui.ErrUnknownView {
		return err
	} else if err == nil {
		footerView.Frame = true
		footerView.Wrap = true
		footerView.Title = ""

		fmt.Fprint(footerView, mainKeys)
	}

	return nil
}

func main() {
	log.Println("Starting LazyDebrid...")
	loadUserSettings()
	getUserTorrents()
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
