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

var userDownloads []DebridDownload
var userSettings = make(map[string]string)
var downloadMap = make(map[string]DebridDownload)
var userConfigPath, _ = os.UserConfigDir()
var lazyDebridConfig = filepath.Join(userConfigPath, "lazyDebrid.json")
var searchQuery string
var views = []string{"search", "torrents", "main"}
var currentViewIdx = 0
var userApiToken string
var userDownloadPath string

func loadUserSettings() error {
	content, err := os.ReadFile(lazyDebridConfig)
	if err != nil {
		fmt.Println("Set API key first.")
	} else {
		_ = json.Unmarshal(content, &userSettings)
	}
	userApiToken = userSettings["apiToken"]
	userDownloadPath = userSettings["downloadPath"]
	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	currentViewIdx = (currentViewIdx + 1) % len(views)
	name := views[currentViewIdx]
	_, err := g.SetCurrentView(name)
	return err
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

func getUserTorrents(url string) map[string]DebridDownload {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
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

	return downloadMap
}

func match(filename, query string) bool {
	return strings.Contains(strings.ToLower(filename), strings.ToLower(query))
}

func renderList(g *gocui.Gui) error {
	v, err := g.View("torrents")
	if err != nil {
		return err
	}
	v.Clear()
	for _, torrentItem := range userDownloads {
		if searchQuery == "" || match(torrentItem.Filename, searchQuery) {
			fmt.Fprintln(v, torrentItem.Filename)
		}
	}
	return nil
}

func searchKeyPress(g *gocui.Gui, v *gocui.View) error {
	searchQuery = strings.TrimSpace(v.Buffer())
	renderList(g)
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

func downloadFile(torrentItem DebridDownload) bool {
	out, err := os.Create(fmt.Sprintf("%s%s", userDownloadPath, torrentItem.Filename))
	defer out.Close()

	if err != nil {
		fmt.Println(err)
	}
	resp, err := http.Get(torrentItem.Download)
	defer resp.Body.Close()
	n, err := io.Copy(out, resp.Body)
	if err == nil {
		fmt.Println(n)
		return true
	}
	return false
}

func generateDetailsString(torrentItem DebridDownload) string {
	detailsString := fmt.Sprintf(
		"ID: %s\nFilename: %s\nMIME Type: %s\nFilesize: %d bytes\nLink: %s\nHost: %s\nHost Icon: %s\nChunks: %d\nDownload: %s\nStreamable: %d\nGenerated: %s\n",
		torrentItem.Id,
		torrentItem.Filename,
		torrentItem.MimeType,
		torrentItem.Filesize,
		torrentItem.Link,
		torrentItem.Host,
		torrentItem.HostIcon,
		torrentItem.Chunks,
		torrentItem.Download,
		torrentItem.Streamable,
		torrentItem.Generated,
	)
	return detailsString
}
func updateDetails(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	title, err := v.Line(cy)
	if err != nil || title == "" {
		return nil
	}

	mainView, err := g.View("main")
	if err != nil {
		return err
	}
	mainView.Clear()

	torrentItem, ok := downloadMap[title]
	if !ok {
		fmt.Fprint(mainView, "No details found.")
		return nil
	}

	detailsString := generateDetailsString(torrentItem)
	fmt.Fprint(mainView, detailsString)
	return nil
}

func selectLine(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return err
	}
	log.Println("Selected:", downloadMap[line].Download)
	return nil
}

func sendLinkToAPI(magnetLink string) error {
	apiURL := "https://api.real-debrid.com/rest/1.0/torrents/addMagnet"
	data := url.Values{}
	data.Set("magnet", magnetLink)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+userApiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		msg, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, msg)
	}

	var result struct {
		ID  string `json:"id"`
		URI string `json:"uri"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	return nil
}

func addMagnetLink(g *gocui.Gui, v *gocui.View) error {
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	maxX, maxY := g.Size()
	if v, err := g.SetView("addMagnet", maxX/4, maxY/4, maxX*3/4, maxY*3/4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Add Magnet Link"
		v.Wrap = true
		v.Editable = true
		v.Clear()
		if _, err := g.SetCurrentView("addMagnet"); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("addMagnet", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		magnetLink := v.Buffer()
		sendLinkToAPI(magnetLink)
		g.DeleteView("addMagnet")
		g.SetCurrentView("torrents")
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func showSearchBar(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("searchPopup", maxX/4, maxY/4, maxX*3/4, maxY/4+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Search"
		v.Editable = true
		v.Wrap = false
		v.Clear()
		g.SetCurrentView("searchPopup")
	}

	if err := g.SetKeybinding("searchPopup", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		searchQuery = strings.TrimSpace(v.Buffer())
		g.DeleteView("searchPopup")
		g.SetCurrentView("torrents")
		return renderList(g)
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("searchPopup", gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		g.DeleteView("searchPopup")
		g.SetCurrentView("torrents")
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("torrents", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("", '/', gocui.ModNone, showSearchBar); err != nil {
		return err
	}
	if err := g.SetKeybinding("torrents", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("torrents", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("torrents", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("torrents", gocui.KeyEnter, gocui.ModNone, downloadSelected); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("search", gocui.KeyEnter, gocui.ModNone, searchKeyPress); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, downloadSelected); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, copyDownloadLink); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlA, gocui.ModNone, addMagnetLink); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, setDownloadPath); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlX, gocui.ModNone, setApiToken); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	splitX := (maxX * 4) / 10
	infoHeight := (maxY - 3) / 4

	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen

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
	}

	if mainView, err := g.SetView("main", splitX+1, 3, maxX-1, maxY-3-infoHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		mainView.Title = "Download Details"
		mainView.Wrap = true
		if len(userDownloads) > 0 {
			first := userDownloads[0]
			mainView.Clear()
			fmt.Fprint(mainView, generateDetailsString(first))
		}
	}

	if infoView, err := g.SetView("info", splitX+1, maxY-3-infoHeight+1, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		infoView.Title = "Info"
		infoView.Wrap = true
		infoView.Autoscroll = true
	}

	if footerView, err := g.SetView("footer", 0, maxY-3, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		footerView.Frame = true
		footerView.Wrap = true
		footerView.Title = "Shortcuts"
		fmt.Fprintln(footerView,
			"TAB: Switch | ↑↓: Navigate | ENTER: Download | Ctrl+A: Add Magnet | Ctrl+C: Copy Link | Ctrl+P: Set Path | Ctrl+X: Set API Key | Ctrl+Q: Quit",
		)
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
func getCurrentLine(v *gocui.View) (string, error) {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	return line, err
}

func main() {
	loadUserSettings()
	getUserTorrents("https://api.real-debrid.com/rest/1.0/downloads?page=1&limit=5000&page=1")
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
