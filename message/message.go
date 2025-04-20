package message

import (
	"fmt"
	"io"
)

type MessageType uint16

const (
	SIGN_UP      MessageType = 0
	SIGN_UP_RESP MessageType = 1
	CONTROL      MessageType = 2
	PLAYER_LOC   MessageType = 3
	PING         MessageType = 4
	PONG         MessageType = 6
	ERROR        MessageType = 7
)

func (mt MessageType) String() string {
	str, ok := map[MessageType]string{
		SIGN_UP:      "SIGN_UP",
		SIGN_UP_RESP: "SIGN_UP_RESP",
		CONTROL:      "CONTROL",
		PLAYER_LOC:   "PLAYER_LOC",
		PING:         "PING",
		PONG:         "PONG",
		ERROR:        "ERROR",
	}[mt]

	if !ok {
		return fmt.Sprintf("INVALID(%v)", uint16(mt))
	}

	return str
}

type PlayerColour struct {
	R uint8
	G uint8
	B uint8
}

func (c PlayerColour) String() string {
	return fmt.Sprintf("RGB(%v, %v, %v)", c.R, c.G, c.B)
}

type Message interface {
	fmt.Stringer
	Read(r io.Reader, msg_len uint32) error
	Write(w io.Writer) error
	Size() uint32
}

type ErrInvalidMessageType MessageType

func (err ErrInvalidMessageType) Error() string {
	return fmt.Sprintf("invalid message type %v", MessageType(err))
}

func ReadMessage(r io.Reader) (Message, error) {
	var zero Message

	h := new(Header)
	if err := h.Read(r); err != nil {
		return zero, err
	}

	var msg Message
	switch h.MsgType {
	case SIGN_UP:
		msg = NewSignUpMessage()
	case SIGN_UP_RESP:
		msg = NewSignUpRespMessage()
	case PING:
		msg = NewPingMessage()
	case PONG:
		msg = NewPongMessage()
	case PLAYER_LOC:
		msg = NewPlayerLocMessage()
	default:
		return zero, ErrInvalidMessageType(h.MsgType)
	}

	if err := msg.Read(r, h.MsgLen); err != nil {
		return zero, err
	}

	return msg, nil

}

func WriteMessage(w io.Writer, msg Message) error {
	h := NewHeader()
	switch msg.(type) {
	case *SignUpMessage:
		h.MsgType = SIGN_UP
	case *SignUpRespMessage:
		h.MsgType = SIGN_UP_RESP
	case *PingMessage:
		h.MsgType = PING
	case *PongMessage:
		h.MsgType = PONG
	case *PlayerLocMessage:
		h.MsgType = PLAYER_LOC
	default:
		return ErrInvalidMessageType(0)
	}

	h.MsgLen = msg.Size()
	err := h.Write(w)
	if err != nil {
		return err
	}

	err = msg.Write(w)
	if err != nil {
		return err
	}

	return nil
}
