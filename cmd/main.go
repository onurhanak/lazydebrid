package main

import (
	"fmt"
	"lazydebrid/internal/actions"
	"lazydebrid/internal/bindings"
	"lazydebrid/internal/config"
	"lazydebrid/internal/views"
	"log"
	"os"
	"path/filepath"

	"github.com/jroimartin/gocui"
)

func init() {
	_, lazyDebridFolderPath, _ := config.ConfigPath()
	logPath := filepath.Join(lazyDebridFolderPath, "lazydebrid.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Could not open log file:", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Println("Starting LazyDebrid...")
	config.HandleFirstRun()

	err := config.LoadUserSettings()
	if err != nil {
		_, lazyDebridFolderPath, _ := config.ConfigPath()
		fmt.Printf("User config is corrupted.\nYou may need to delete %s folder manually.\n", lazyDebridFolderPath)
		log.Fatalf("Could not load user settings, bailing. Error: %s", err)
	}
	actions.GetUserTorrents()
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetManagerFunc(views.Layout)

	if err := bindings.Keybindings(g); err != nil {
		log.Panicln(err)
	}

	// delay populate views until views are ready
	// otherwise active torrents does not show
	views.OnLayoutReady = views.PopulateViews
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("GUI loop error: %v", err)
	}

}
