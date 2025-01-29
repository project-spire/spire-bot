package core

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	Bots uint

	AuthHost string
	AuthPort int

	GameHost string
	GamePort int
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
