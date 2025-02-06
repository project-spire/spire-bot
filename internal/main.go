package main

import (
	"net"
	"strconv"
	"sync"

	"spire/bot/internal/bot"
	"spire/bot/internal/core"
)

func main() {
	settings := core.ReadSettings("settings.yaml")
	wg := sync.WaitGroup{}
	wg.Add(settings.Bots)

	launchBot := func(id uint64) {
		defer wg.Done()

		b := bot.NewBot(id)
		go b.Start(
			net.JoinHostPort(settings.LobbyHost, strconv.Itoa(settings.LobbyPort)),
			net.JoinHostPort(settings.GameHost, strconv.Itoa(settings.GamePort)))

		<-b.Stopped
	}

	for i := 1; i <= settings.Bots; i++ {
		go launchBot(uint64(i))
	}

	wg.Wait()
}
