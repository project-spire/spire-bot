package core

import (
	"encoding/binary"
	"io"
	"log/slog"
	"net"
	"strconv"

	"google.golang.org/protobuf/proto"
	"spire/bot/gen/msg"
)

type Connection struct {
	Receiver chan *msg.BaseMessage
	Sender   chan *msg.BaseMessage
	conn     net.Conn
	stop     chan struct{}
}

func NewConnection(host string, port int) (*Connection, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}

	return &Connection{
		Receiver: make(chan *msg.BaseMessage),
		Sender:   make(chan *msg.BaseMessage),
		conn:     conn,
		stop:     make(chan struct{}),
	}, nil
}

func (c *Connection) Start() {
	go c.receive()
	go c.send()
}

func (c *Connection) Stop() {
	close(c.stop)
	close(c.Receiver)
	close(c.Sender)
	c.conn.Close()
}

func (c *Connection) receive() {
	for {
		select {
		case <-c.stop:
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
				slog.Warn("Error unmarshal InMessage")
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
		case <-c.stop:
			return
		case message := <-c.Sender:
			buf, err := proto.MarshalOptions{}.MarshalAppend(make([]byte, 2), message)
			if err != nil {
				c.Stop()
				return
			}

			if 4+len(buf) > 65536 {
				slog.Warn("OutMessage to large: %v", len(buf))
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
