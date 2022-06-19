package main

import (
	"encoding/json"
	"io/ioutil"
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
	data, _ := json.MarshalIndent(config, "", "\t")
	ioutil.WriteFile("config.json", data, 0644)
}

func LoadConfig() Config {
	config := Config{}

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return config
	}

	json.Unmarshal([]byte(data), &config)

	return config
}
