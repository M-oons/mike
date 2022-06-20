package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	Hotkeys []ConfigHotkey `json:"hotkeys"`
}

type ConfigHotkey struct {
	Action   string `json:"action"`
	KeyCode  int    `json:"keycode"`
	Ctrl     bool   `json:"ctrl"`
	Shift    bool   `json:"shift"`
	Alt      bool   `json:"alt"`
	Win      bool   `json:"win"`
	NoRepeat bool   `json:"norepeat"`
}

func (config Config) Save() {
	if !ensureConfig() {
		return
	}

	writeConfig(config)
}

func LoadConfig() Config {
	config := Config{}

	if !ensureConfig() {
		return config
	}

	dir := getConfigPath()
	if dir == "" {
		return config
	}

	data := readConfig()

	json.Unmarshal(data, &config)

	return config
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
		config := Config{
			Hotkeys: make([]ConfigHotkey, 0),
		}
		return writeConfig(config)
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

	return filepath.Join(roaming, AppAuthor, AppName)
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
