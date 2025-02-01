package core

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	Bots     int    `yaml:"bots"`
	AuthHost string `yaml:"auth_host"`
	AuthPort int    `yaml:"auth_port"`
	GameHost string `yaml:"game_host"`
	GamePort int    `yaml:"game_port"`
}

func ReadSettings(settingsPath string) Settings {
	settings := Settings{}

	settingsData, err := os.ReadFile(settingsPath)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", settingsPath, err)
	}

	err = yaml.Unmarshal(settingsData, &settings)
	if err != nil {
		log.Fatalf("Failed to parse %s: %v", settingsPath, err)
	}

	return settings
}
