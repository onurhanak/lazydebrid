package main

import (
	"lazydebrid/internal/actions"
	"lazydebrid/internal/bindings"
	"lazydebrid/internal/config"
	"lazydebrid/internal/handlers"
	"lazydebrid/internal/views"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

func init() {
	logFile, err := os.OpenFile("lazydebrid.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Could not open log file:", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Println("Starting LazyDebrid...")
	err := config.LoadUserSettings()
	if err != nil {
		log.Fatalln("Could not not load user settings, bailing.")
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

	g.Update(func(g *gocui.Gui) error {
		handlers.PopulateViews(g)
		v, err := g.View("torrents")
		if err != nil {
			return err
		}
		return handlers.UpdateDetails(g, v)
	})

	// delay populate views until views are ready
	// otherwise active torrents does not show
	views.OnLayoutReady = handlers.PopulateViews
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}
