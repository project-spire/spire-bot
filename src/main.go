//go:generate cp settings.yaml ./build/settings.yaml

package main

import (
	"log/slog"
	"net"
	"strconv"

	"spire/bot/src/core"
)

func main() {
	settings := core.ReadSettings("settings.yaml")

	conn, err := net.Dial("tcp", net.JoinHostPort(settings.GameHost, strconv.Itoa(settings.GamePort)))
	if err != nil {
		slog.Error("Failed to connect to %s:%d. %v", settings.GameHost, settings.GamePort, err)
	}
	defer conn.Close()

	slog.Info("Successfully connected to %s:%d", settings.GameHost, settings.GamePort)
}
