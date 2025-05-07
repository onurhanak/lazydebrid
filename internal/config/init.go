package config

import (
	"bufio"
	"fmt"
	"lazydebrid/internal/logs"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func HandleFirstRun() {
	if isFirstRun() {
		if err := SetupConfigFromUserInput(); err != nil {
			logs.LogEvent(fmt.Errorf("failed to set up config: %w", err))
			fmt.Println("Failed to save config:", err)
			os.Exit(1)
		}
	}
}

func isFirstRun() bool {
	configPath, dirPath, err := ConfigPath()
	if err != nil {
		logs.LogEvent(fmt.Errorf("cannot determine config path: %w", err))
		return false
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			logs.LogEvent(fmt.Errorf("cannot create config folder: %w", err))
			return false
		}
		return true
	}
	return false
}

func SetupConfigFromUserInput() error {
	fmt.Println("Detected first run.")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your API token: ")
	apiToken, _ := reader.ReadString('\n')
	apiToken = strings.TrimSpace(apiToken)
	if apiToken == "" {
		return fmt.Errorf("API token cannot be empty")
	}

	fmt.Print("Enter download path (leave empty for default $HOME/Downloads): ")
	downloadPath, _ := reader.ReadString('\n')
	downloadPath = strings.TrimSpace(downloadPath)

	if downloadPath == "" {
		usr, err := user.Current()
		if err != nil {
			return fmt.Errorf("cannot get current user: %w", err)
		}
		downloadPath = filepath.Join(usr.HomeDir, "Downloads")
	}

	if stat, err := os.Stat(downloadPath); os.IsNotExist(err) {
		fmt.Printf("Download path '%s' does not exist. Create it? (y/n): ", downloadPath)
		resp, _ := reader.ReadString('\n')
		resp = strings.TrimSpace(resp)
		if strings.ToLower(resp) == "y" {
			if err := os.MkdirAll(downloadPath, 0755); err != nil {
				return fmt.Errorf("failed to create download directory: %w", err)
			}
		} else {
			return fmt.Errorf("download path must exist; exiting")
		}
	} else if !stat.IsDir() {
		return fmt.Errorf("download path '%s' is not a directory", downloadPath)
	}

	// Save settings
	if err := SaveSetting("apiToken", apiToken); err != nil {
		return fmt.Errorf("failed to save apiToken: %w", err)
	}
	if err := SaveSetting("downloadPath", downloadPath); err != nil {
		return fmt.Errorf("failed to save downloadPath: %w", err)
	}

	fmt.Println("Configuration saved.")
	return nil
}
