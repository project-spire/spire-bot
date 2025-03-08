package core

import (
	"encoding/binary"
	"io"
	"log/slog"
	"net"
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

func (c *Connection) ConnectAsync(address string) <-chan error {
	errResult := make(chan error, 1)

	go func() {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			c.logger.Warn("Failed to connect to %s: %v", address, err)
			errResult <- err
			close(errResult)
			return
		}

		c.conn = conn

		close(errResult)
	}()

	return errResult
}

func (c *Connection) Start(address string) {
	go c.receive()
	go c.send()
}

func (c *Connection) Stop() {
	c.stopOnce.Do(func() {
		if c.conn != nil {
			c.conn.Close()
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
			return

		default:
			headerBuf := make([]byte, 2)
			if _, err := io.ReadFull(c.conn, headerBuf); err != nil {
				c.logger.Error("Error receiving header: %v", err)
				c.Stop()
				return
			}

			bodyLen := binary.BigEndian.Uint16(headerBuf)
			bodyBuf := make([]byte, bodyLen)
			if _, err := io.ReadFull(c.conn, bodyBuf); err != nil {
				c.logger.Error("Error receiving body: %v", err)
				c.Stop()
				return
			}

			//base := &msg.BaseMessage{}
			//if err := proto.Unmarshal(bodyBuf, base); err != nil {
			//	c.logger.Warn("Error unmarshal InMessage")
			//	c.Stop()
			//	return
			//}
			c.Receiver <- bodyBuf
		}
	}
}

func (c *Connection) send() {
	for {
		select {
		case <-c.Stopped:
			return

		case buf := <-c.Sender:
			if _, err := c.conn.Write(buf); err != nil {
				c.logger.Error("Error sending: %v", err)
				c.Stop()
				return
			}
		}
	}
}
