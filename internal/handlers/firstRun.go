package handlers

import (
	"fmt"
	"lazydebrid/internal/config"
)

func HandleFirstRun() {
	if config.CheckFirstRun() {
		if err := config.SetupConfigFromUserInput(); err != nil {
			fmt.Println("Failed to save config:", err)
		}
	}
}
