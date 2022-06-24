package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/m-oons/mike/info"
)

var Current Config

type Config struct {
	Hotkeys []Hotkey `json:"hotkeys"`
	Sounds  bool     `json:"sounds"`
}

type Hotkey struct {
	Action   string `json:"action"`
	Key      string `json:"key"`
	Ctrl     bool   `json:"ctrl"`
	Shift    bool   `json:"shift"`
	Alt      bool   `json:"alt"`
	Win      bool   `json:"win"`
	NoRepeat bool   `json:"norepeat"`
}

func SaveConfig() {
	if !ensureConfig() {
		return
	}

	writeConfig(Current)
}

func LoadConfig() {
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

	return filepath.Join(roaming, info.Author, info.AppName)
}

func readConfig() []byte {
	dir := getConfigPath()
	if dir == "" {
		return nil
	}

	data, err := ioutil.ReadFile(filepath.Join(dir, "config.json"))
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
	err := ioutil.WriteFile(filepath.Join(dir, "config.json"), data, 0644)

	return err == nil
}

func defaultConfig() Config {
	return Config{
		Hotkeys: make([]Hotkey, 0),
		Sounds:  true,
	}
}
