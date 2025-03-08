package bot

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"spire/bot/internal/core"
)

type Account struct {
	AccountId   uint64
	CharacterId uint64
	AuthToken   string
}

type Bot struct {
	BotId   uint64
	Stopped chan struct{}
	Account Account

	conn     *core.Connection
	logger   *slog.Logger
	stopOnce sync.Once
}

func NewBot(id uint64) *Bot {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("id", id)

	return &Bot{
		BotId:    id,
		Stopped:  make(chan struct{}, 1),
		conn:     core.NewConnection(logger),
		logger:   logger,
		stopOnce: sync.Once{},
	}
}

func (b *Bot) Start(lobbyAddress string, gameAddress string) {
	if b.RequestAccount(lobbyAddress) != nil {
		b.Stop()
		return
	}

	connErr := b.conn.ConnectAsync(gameAddress)
	if err := <-connErr; err != nil {
		b.Stop()
		return
	}
	b.logger.Info(fmt.Sprintf("Connected to %s", gameAddress))

	b.conn.Start()
	b.RequestLogin()

	go func() {
		//<-b.conn.Stopped
		//b.Stop()
	}()
}

func (b *Bot) Stop() {
	b.stopOnce.Do(func() {
		b.logger.Info("Bot Stopped")
		b.conn.Stop()
		close(b.Stopped)
	})
}
