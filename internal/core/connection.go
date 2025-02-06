package core

import (
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
	"log/slog"
	"net"
	"spire/bot/gen/msg"
	"sync"
)

type Connection struct {
	Receiver chan *msg.BaseMessage
	Sender   chan *msg.BaseMessage
	Stopped  chan struct{}

	conn     net.Conn
	logger   *slog.Logger
	stopOnce sync.Once
}

func NewConnection(logger *slog.Logger) *Connection {
	return &Connection{
		Receiver: make(chan *msg.BaseMessage),
		Sender:   make(chan *msg.BaseMessage),
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
				c.Stop()
				return
			}

			bodyLen := binary.BigEndian.Uint16(headerBuf)
			bodyBuf := make([]byte, bodyLen)
			if _, err := io.ReadFull(c.conn, bodyBuf); err != nil {
				c.Stop()
				return
			}

			base := &msg.BaseMessage{}
			if err := proto.Unmarshal(bodyBuf, base); err != nil {
				c.logger.Warn("Error unmarshal InMessage")
				c.Stop()
				return
			}
			c.Receiver <- base
		}
	}
}

func (c *Connection) send() {
	for {
		select {
		case <-c.Stopped:
			return

		case message := <-c.Sender:
			buf, err := proto.MarshalOptions{}.MarshalAppend(make([]byte, 2), message)
			if err != nil {
				c.Stop()
				return
			}

			if 2+len(buf) > 65536 {
				c.logger.Warn("OutMessage to large: %v", len(buf))
				c.Stop()
				return
			}
			binary.BigEndian.PutUint16(buf[:2], uint16(len(buf)-2))

			if _, err := c.conn.Write(buf); err != nil {
				c.Stop()
				return
			}
		}
	}
}
