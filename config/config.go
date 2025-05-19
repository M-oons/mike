package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/m-oons/mike/info"
)

var Current Config

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
	Voicemeeter ConfigControllerVoicemeeter `json:"voicemeeter"`
}

type ConfigControllerVoicemeeter struct {
	RemoteDLLPath string `json:"remoteDLLPath"`
	Output        byte   `json:"output"`
}

func Save() {
	if !ensureConfig() {
		return
	}

	writeConfig(Current)
}

func Load() {
	Current = Config{}

	if !ensureConfig() {
		return
	}

	dir := getConfigPath()
	if dir == "" {
		return
	}

	data := readConfig()
	json.Unmarshal(data, &Current)
}

func ensureConfig() bool {
	dir := getConfigPath()
	if dir == "" {
		return false
	}

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return false
	}

	if !configFileExists() {
		return writeConfig(defaultConfig())
	}

	return true
}

func configFileExists() bool {
	dir := getConfigPath()
	if dir == "" {
		return false
	}

	if _, err := os.Stat(filepath.Join(dir, "config.json")); err != nil {
		return false
	}

	return true
}

func getConfigPath() string {
	roaming, err := os.UserConfigDir()
	if err != nil {
		return ""
	}

	return filepath.Join(roaming, info.AppName)
}

func readConfig() []byte {
	dir := getConfigPath()
	if dir == "" {
		return nil
	}

	data, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		return nil
	}

	return data
}

func writeConfig(config Config) bool {
	dir := getConfigPath()
	if dir == "" {
		return false
	}

	data, _ := json.MarshalIndent(config, "", "\t")
	err := os.WriteFile(filepath.Join(dir, "config.json"), data, 0644)

	return err == nil
}

func defaultConfig() Config {
	return Config{
		Hotkeys: make([]ConfigHotkey, 0),
		Sounds: ConfigSounds{
			Enabled: true,
			Volume:  100,
		},
		Controller: ConfigController{
			Type: "windows",
			Voicemeeter: ConfigControllerVoicemeeter{
				RemoteDLLPath: "C:/Program Files (x86)/VB/Voicemeeter/VoicemeeterRemote64.dll",
				Output:        1,
			},
		},
	}
}
