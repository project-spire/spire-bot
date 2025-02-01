//go:generate cp settings.yaml ./build/settings.yaml

package main

import (
	"log/slog"
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

	launchBot := func(id int) {
		defer wg.Done()

		b := bot.NewBot(id)
		err := <-b.StartAsync(net.JoinHostPort(settings.GameHost, strconv.Itoa(settings.GamePort)))
		if err != nil {
			slog.Error("Bot failed to start: %v", err)
			return
		}

		<-b.Stopped
	}

	for i := 1; i <= settings.Bots; i++ {
		go launchBot(i)
	}

	wg.Wait()
}
