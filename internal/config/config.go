package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const configFile = "lazyDebrid.json"

var (
	UserSettings = make(map[string]string)

	UserApiToken     string
	UserDownloadPath string

	UserConfigPath, _ = os.UserConfigDir()
	LazyDebridConfig  = filepath.Join(UserConfigPath, "lazyDebrid.json")
	SearchQuery       string
)

func ConfigPath() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, configFile)
}

func LoadUserSettings() error {
	data, err := os.ReadFile(ConfigPath())
	if err == nil {
		_ = json.Unmarshal(data, &UserSettings)
	}
	UserApiToken = strings.TrimSpace(UserSettings["apiToken"])
	UserDownloadPath = strings.TrimSpace(UserSettings["downloadPath"])
	return nil
}

func SaveSetting(key, value string) error {
	value = strings.TrimSpace(value)
	data, _ := os.ReadFile(ConfigPath())
	_ = json.Unmarshal(data, &UserSettings)
	UserSettings[key] = value
	content, err := json.MarshalIndent(UserSettings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(), content, 0644)
}
