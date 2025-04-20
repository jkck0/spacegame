package message

import (
	"io"
)

type PingMessage struct{}

func NewPingMessage() *PingMessage {
	return &PingMessage{}
}

func (m *PingMessage) Read(r io.Reader, _ uint32) error {
	return nil
}

func (m *PingMessage) Write(w io.Writer) error {
	return nil
}

func (m *PingMessage) Size() uint32 {
	return 0
}

func (m *PingMessage) String() string {
	return "Ping"
}

type PongMessage struct{}

func NewPongMessage() *PongMessage {
	return &PongMessage{}
}

func (m *PongMessage) Read(r io.Reader, _ uint32) error {
	return nil
}

func (m *PongMessage) Write(w io.Writer) error {
	return nil
}

func (m *PongMessage) Size() uint32 {
	return 0
}

func (m *PongMessage) String() string {
	return "Pong"
}
