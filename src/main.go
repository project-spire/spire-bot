//go:generate cp settings.yaml ./build/settings.yaml

package main

import (
	"log"
	"net"
	"strconv"

	"spire/bot/src/core"
)

func main() {
	settings := core.ReadSettings("settings.yaml")

	conn, err := net.Dial("tcp", net.JoinHostPort(settings.GameHost, strconv.Itoa(settings.GamePort)))
	if err != nil {
		log.Fatalf("Failed to connect to %s:%d. %v", settings.GameHost, settings.GamePort, err)
	}
	defer conn.Close()

	log.Printf("Successfully connected to %s:%d", settings.GameHost, settings.GamePort)
}
