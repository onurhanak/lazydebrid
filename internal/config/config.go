package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	configFileName      = "lazyDebrid.json"
	lazyDebridFolder    = "lazyDebrid"
	defaultAPIToken     = ""
	defaultDownloadPath = ""
)

var (
	settings         = make(map[string]string)
	userApiToken     string
	userDownloadPath string
	searchQuery      string
)

func ConfigPath() (string, string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", "", fmt.Errorf("failed to get user config dir: %w", err)
	}
	lazyDebridFolderPath := filepath.Join(configDir, lazyDebridFolder)
	lazyDebridConfigPath := filepath.Join(lazyDebridFolderPath, configFileName)
	return lazyDebridConfigPath, lazyDebridFolderPath, nil
}

func EnsureConfigExists() error {
	configPath, dirPath, err := ConfigPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create config dir: %w", err)
		}
		defaultSettings := map[string]string{
			"apiToken":     defaultAPIToken,
			"downloadPath": defaultDownloadPath,
		}
		content, err := json.MarshalIndent(defaultSettings, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to serialize default settings: %w", err)
		}
		if err := os.WriteFile(configPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write default config: %w", err)
		}
	}
	return nil
}

func LoadUserSettings() error {
	if err := EnsureConfigExists(); err != nil {
		return err
	}

	configPath, _, err := ConfigPath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("failed to parse config JSON: %w", err)
	}

	userApiToken = strings.TrimSpace(settings["apiToken"])
	userDownloadPath = strings.TrimSpace(settings["downloadPath"])
	return nil
}

func SaveSetting(key, value string) error {
	configPath, _, err := ConfigPath()
	if err != nil {
		return err
	}

	// Re-load the settings
	data, err := os.ReadFile(configPath)
	if err == nil {
		if err := json.Unmarshal(data, &settings); err != nil {
			return fmt.Errorf("failed to reload existing config: %w", err)
		}
	}

	settings[key] = strings.TrimSpace(value)

	content, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	if err := os.WriteFile(configPath, content, 0644); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Sync globals if key was updated
	switch key {
	case "apiToken":
		userApiToken = settings[key]
	case "downloadPath":
		userDownloadPath = settings[key]
	}

	return nil
}

func Get(key string) (string, error) {
	val, ok := settings[key]
	if !ok {
		return "", fmt.Errorf("setting not found: %s", key)
	}
	return val, nil
}

func APIToken() string     { return userApiToken }
func DownloadPath() string { return userDownloadPath }
func SearchQuery() string  { return searchQuery }
func SetSearchQuery(q string) {
	searchQuery = strings.TrimSpace(q)
}
