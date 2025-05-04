package utils

import (
	"fmt"
	"lazydebrid/internal/actions"
	"lazydebrid/internal/config"
	"lazydebrid/internal/models"
	"strings"

	"github.com/jroimartin/gocui"
)

func Match(filename, query string) bool {
	return strings.Contains(strings.ToLower(filename), strings.ToLower(query))
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

	for _, torrentItem := range actions.UserDownloads {
		if config.SearchQuery() == "" || Match(torrentItem.Filename, config.SearchQuery()) {
			_, err := fmt.Fprintln(v, torrentItem.Filename)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetCurrentLine(v *gocui.View) (string, error) {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	return line, err
}

func GenerateDetailsString(torrentItem models.Torrent) string {
	detailsString := fmt.Sprintf(
		"ID: %s\nFilename: %s\nFilesize: %d bytes\nLink: %s\nDownload: %s\nStreamable: %d",
		torrentItem.ID,
		torrentItem.Filename,
		torrentItem.Bytes,
		torrentItem.Hash,
		torrentItem.Hash,
		torrentItem.Bytes,
	)
	return detailsString
}
