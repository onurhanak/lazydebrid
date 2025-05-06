package config

import (
	"bufio"
	"fmt"
	"lazydebrid/internal/logs"
	"os"
	"os/user"
	"path/filepath"
)

func HandleFirstRun() {
	if CheckFirstRun() {
		if err := SetupConfigFromUserInput(); err != nil {
			fmt.Println("Failed to save config:", err)
		}
	}
}

func CheckFirstRun() bool {
	lazyDebridConfigPath, lazyDebridFolderPath, err := ConfigPath()
	if err != nil {
		logs.LogEvent(fmt.Errorf("Cannot read config: %s", err))
		return false
	}

	if _, err := os.Stat(lazyDebridConfigPath); os.IsNotExist(err) {
		if err := os.MkdirAll(lazyDebridFolderPath, os.ModePerm); err != nil {

			logs.LogEvent(fmt.Errorf("Cannot create config folder: %s", err))
			return false
		}

		file, err := os.Create(lazyDebridConfigPath)
		if err != nil {
			logs.LogEvent(fmt.Errorf("Cannot create config file: %s", err))
			return false
		}
		defer file.Close()

		return true // first run
	}

	return false
}

func SetupConfigFromUserInput() error {
	fmt.Println("Detected first run.")
	fmt.Print("Enter your API token: ")

	reader := bufio.NewReader(os.Stdin)
	apiToken, _ := reader.ReadString('\n')
	apiToken = trimNewline(apiToken)
	if len(apiToken) <= 0 {
		return fmt.Errorf("API token cannot be empty.")
	}

	fmt.Print("Enter download path (leave empty for default $HOME/Downloads): ")
	downloadPath, _ := reader.ReadString('\n')
	downloadPath = trimNewline(downloadPath)

	if downloadPath == "" {
		usr, _ := user.Current()
		downloadPath = filepath.Join(usr.HomeDir, "Downloads/")
	} else {
		if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
			fmt.Printf("Download path '%s' does not exist. Create it? (y/n): ", downloadPath)
			var response string
			fmt.Scanln(&response)
			if response == "y" || response == "Y" {
				err := os.MkdirAll(downloadPath, 0755)
				if err != nil {
					fmt.Println("Failed to create directory:", err)
					os.Exit(1)
				}
			} else {
				// delete config folder so it triggers first run next time
				_, lazyDebridFolderPath, err := ConfigPath()
				if err == nil {
					_ = os.RemoveAll(lazyDebridFolderPath)
				}
				fmt.Println("Downloads folder must exist. Exiting.")
				os.Exit(1)
			}
		}
	}

	if err := SaveSetting("apiToken", apiToken); err != nil {
		return err
	}
	if err := SaveSetting("downloadPath", downloadPath); err != nil {
		return err
	}

	fmt.Println("Configuration saved.")
	return nil
}

func trimNewline(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == '\r' {
		s = s[:len(s)-1]
	}
	return s
}
