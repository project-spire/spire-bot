package core

import (
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"net"
	"spire/protocol"
	"sync"
)

type Connection struct {
	Receiver chan []byte
	Sender   chan []byte
	Stopped  chan struct{}

	conn     net.Conn
	logger   *slog.Logger
	stopOnce sync.Once
}

func NewConnection(logger *slog.Logger) *Connection {
	return &Connection{
		Receiver: make(chan []byte, 16),
		Sender:   make(chan []byte, 16),
		Stopped:  make(chan struct{}, 1),
		conn:     nil,
		logger:   logger,
		stopOnce: sync.Once{},
	}
}

func (c *Connection) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		c.logger.Warn("Failed to connect to %s: %v", address, err)
		return err
	}

	c.conn = conn
	return nil
}

func (c *Connection) Start() {
	go c.receive()
	go c.send()
}

func (c *Connection) Stop() {
	c.stopOnce.Do(func() {
		c.logger.Info("Closing connection")
		if c.conn != nil {
			_ = c.conn.Close()
		}

		close(c.Receiver)
		close(c.Sender)
		close(c.Stopped)
	})
}

func (c *Connection) receive() {
	for {
		select {
		case <-c.Stopped:
			c.logger.Info("receive stopped")
			return

		default:
			headerBuf := make([]byte, protocol.HeaderSize)
			if _, err := io.ReadFull(c.conn, headerBuf); err != nil {
				c.logger.Error(fmt.Sprintf("Error receiving header: %v", err))
				c.Stop()
				return
			}

			bodyLen := binary.BigEndian.Uint16(headerBuf)
			bodyBuf := make([]byte, bodyLen)
			if _, err := io.ReadFull(c.conn, bodyBuf); err != nil {
				c.logger.Error(fmt.Sprintf("Error receiving body: %v", err))
				c.Stop()
				return
			}

			c.Receiver <- bodyBuf
		}
	}
}

func (c *Connection) send() {
	for {
		select {
		case <-c.Stopped:
			c.logger.Info("send stopped")
			return

		case buf := <-c.Sender:
			if _, err := c.conn.Write(buf); err != nil {
				c.logger.Error(fmt.Sprintf("Error sending: %v", err))
				c.Stop()
				return
			}
		}
	}
}
