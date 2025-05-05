package config

import (
	"bufio"
	"fmt"
	"lazydebrid/internal/logs"
	"os"
	"os/user"
	"path/filepath"
)

func CheckFirstRun() bool {
	lazyDebridConfigPath, lazyDebridFolderPath, err := ConfigPath()
	if err != nil {
		logs.LogEvent(err)
		return false
	}

	if _, err := os.Stat(lazyDebridConfigPath); os.IsNotExist(err) {
		if err := os.MkdirAll(lazyDebridFolderPath, os.ModePerm); err != nil {
			logs.LogEvent(err)
			return false
		}

		file, err := os.Create(lazyDebridConfigPath)
		if err != nil {
			logs.LogEvent(err)
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

	fmt.Print("Enter download path (leave empty for default $HOME/Downloads): ")
	downloadPath, _ := reader.ReadString('\n')
	downloadPath = trimNewline(downloadPath)

	if downloadPath == "" {
		usr, _ := user.Current()
		downloadPath = filepath.Join(usr.HomeDir, "Downloads/")
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
