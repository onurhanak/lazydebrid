package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const configFileName = "lazyDebrid.json"
const lazyDebridFolder = "lazyDebrid"

var (
	settings = make(map[string]string)

	userApiToken     string
	userDownloadPath string
	searchQuery      string
)

func ConfigPath() (string, string, error) {
	configDir, err := os.UserConfigDir()
	lazyDebridFolderPath := filepath.Join(configDir, lazyDebridFolder)
	lazyDebridConfigPath := filepath.Join(lazyDebridFolderPath, configFileName)
	if err != nil {
		return "", "", err
	}

	return lazyDebridConfigPath, lazyDebridFolderPath, nil
}

func LoadUserSettings() error {
	lazyDebridConfigPath, _, err := ConfigPath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(lazyDebridConfigPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &settings); err != nil {
		return err
	}

	userApiToken = strings.TrimSpace(settings["apiToken"])
	userDownloadPath = strings.TrimSpace(settings["downloadPath"])
	return nil
}

func SaveSetting(key, value string) error {
	lazyDebridConfigPath, _, err := ConfigPath()
	if err != nil {
		return err
	}
	// reload in case there is manual modification
	data, _ := os.ReadFile(lazyDebridConfigPath)
	_ = json.Unmarshal(data, &settings)

	settings[key] = strings.TrimSpace(value)

	content, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(lazyDebridConfigPath, content, 0644)
}

func Get(key string) (string, error) {
	val, ok := settings[key]
	if !ok {
		return "", errors.New("setting not found")
	}
	return val, nil
}

func APIToken() string     { return userApiToken }
func DownloadPath() string { return userDownloadPath }
func SearchQuery() string  { return searchQuery }
func SetSearchQuery(q string) {
	searchQuery = strings.TrimSpace(q)
}
