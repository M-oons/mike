package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/m-oons/mike/internal/info"
)

type Config struct {
	Hotkeys    []ConfigHotkey   `json:"hotkeys"`
	Sounds     ConfigSounds     `json:"sounds"`
	Controller ConfigController `json:"controller"`
}

type ConfigHotkey struct {
	Action   string `json:"action"`
	Key      string `json:"key"`
	Ctrl     bool   `json:"ctrl"`
	Shift    bool   `json:"shift"`
	Alt      bool   `json:"alt"`
	Win      bool   `json:"win"`
	NoRepeat bool   `json:"norepeat"`
}

type ConfigSounds struct {
	Enabled bool `json:"enabled"`
	Volume  int  `json:"volume"`
}

type ConfigController struct {
	Type        string                      `json:"type"`
	Windows     ConfigControllerWindows     `json:"windows"`
	Voicemeeter ConfigControllerVoicemeeter `json:"voicemeeter"`
}

type ConfigControllerWindows struct{}

type ConfigControllerVoicemeeter struct {
	RemoteDLLPath string `json:"remoteDLLPath"`
	Parameter     string `json:"parameter"`
}

func Load() (*Config, error) {
	config := Config{}

	if err := ensureConfig(); err != nil {
		return nil, fmt.Errorf("error ensuring config exists: %w", err)
	}

	data, err := readConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file")
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config JSON")
	}

	return &config, nil
}

func ensureConfig() error {
	dir, err := getConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	if !configFileExists() {
		return writeConfig(defaultConfig())
	}

	return nil
}

func configFileExists() bool {
	dir, err := getConfigPath()
	if err != nil {
		return false
	}

	if _, err := os.Stat(filepath.Join(dir, "config.json")); err != nil {
		return false
	}

	return true
}

func getConfigPath() (string, error) {
	roaming, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(roaming, strings.ToLower(info.AppName)), nil
}

func readConfig() ([]byte, error) {
	dir, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func writeConfig(config *Config) error {
	dir, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(&config, "", "	")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, "config.json"), data, 0644)
}

func defaultConfig() *Config {
	return &Config{
		Hotkeys: make([]ConfigHotkey, 0),
		Sounds: ConfigSounds{
			Enabled: true,
			Volume:  100,
		},
		Controller: ConfigController{
			Type:    "windows",
			Windows: ConfigControllerWindows{},
			Voicemeeter: ConfigControllerVoicemeeter{
				RemoteDLLPath: "C:/Program Files (x86)/VB/Voicemeeter/VoicemeeterRemote64.dll",
				Parameter:     "Bus[2]",
			},
		},
	}
}
