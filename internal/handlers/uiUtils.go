package handlers

import (
	"fmt"
	"lazydebrid/internal/actions"
	"lazydebrid/internal/models"

	"github.com/jroimartin/gocui"
)

func getSelectedItem(v *gocui.View) (models.Download, error) {
	_, cy := v.Cursor()
	line, err := v.Line(cy)
	if err != nil {
		return models.Download{}, fmt.Errorf("unable to get selected line: %w", err)
	}
	item, ok := actions.FilesMap[line]
	if !ok {
		return models.Download{}, fmt.Errorf("no download item found for selected line")
	}
	return item, nil
}
