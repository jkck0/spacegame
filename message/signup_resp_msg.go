package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

type SignUpRespMessage struct {
	ID      uint32
	MapSize uint16
}

func NewSignUpRespMessage() *SignUpRespMessage {
	return &SignUpRespMessage{}
}

func (m *SignUpRespMessage) Read(r io.Reader, _ uint32) error {
	err := binary.Read(r, binary.BigEndian, m)
	if err != nil {
		return err
	}

	return nil
}

func (m *SignUpRespMessage) Write(w io.Writer) error {
	err := binary.Write(w, binary.BigEndian, m)
	if err != nil {
		return err
	}

	return nil
}

func (m *SignUpRespMessage) Size() uint32 {
	return 6
}

var signUpRespFormat string = `SignUpRespMessage{
    ID:       %v
    MapSize:  %v
}`

func (m *SignUpRespMessage) String() string {
	return fmt.Sprintf(signUpRespFormat, m.ID, m.MapSize)
}
