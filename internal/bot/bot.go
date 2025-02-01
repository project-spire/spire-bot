package bot

import (
	"fmt"
	"log/slog"
	"os"
	"spire/bot/internal/core"
	"sync"
)

type Bot struct {
	Id      int
	Stopped chan struct{}

	conn     *core.Connection
	logger   *slog.Logger
	stopOnce sync.Once
}

func NewBot(id int) *Bot {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("id", id)

	return &Bot{
		Id:       id,
		Stopped:  make(chan struct{}, 1),
		conn:     core.NewConnection(logger),
		logger:   logger,
		stopOnce: sync.Once{},
	}
}

func (b *Bot) StartAsync(gameAddress string) <-chan error {
	errResult := make(chan error, 1)

	go func() {
		connErr := b.conn.ConnectAsync(gameAddress)
		if err := <-connErr; err != nil {
			errResult <- err
			close(errResult)
			return
		}
		b.logger.Info(fmt.Sprintf("Connected to %s", gameAddress))
		//slog.Info(fmt.Sprintf("%s connected to %s", b.Display(), gameAddress))

		b.conn.Start(gameAddress)

		go func() {
			<-b.conn.Stopped
			b.Stop()
		}()

		close(errResult)
	}()

	return errResult
}

func (b *Bot) Stop() {
	b.stopOnce.Do(func() {
		b.logger.Info("Stopped")
		b.conn.Stop()
		close(b.Stopped)
	})
}

func (b *Bot) Display() string {
	return fmt.Sprintf("Bot { id: %d }", b.Id)
}
