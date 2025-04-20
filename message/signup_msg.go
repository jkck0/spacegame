package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

type SignUpMessage struct {
	PlayerColour
	NameLen uint16
	Name    []byte
}

func NewSignUpMessage() *SignUpMessage {
	return &SignUpMessage{}
}

func (m *SignUpMessage) Read(r io.Reader, _ uint32) error {
	err := binary.Read(r, binary.BigEndian, m.PlayerColour)
	if err != nil {
		return err
	}

	var buf [3]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return err
	}

	m.NameLen = binary.BigEndian.Uint16(buf[1:])
	m.Name = make([]byte, m.NameLen)

	if _, err := io.ReadFull(r, m.Name); err != nil {
		return err
	}

	return nil
}

func (m *SignUpMessage) Write(w io.Writer) error {
	if len(m.Name) > int(m.NameLen) {
		return fmt.Errorf("name is longer than Message.NameLen")
	}

	err := binary.Write(w, binary.BigEndian, m.PlayerColour)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte{0})
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.BigEndian, m.NameLen)
	if err != nil {
		return err
	}

	_, err = w.Write(m.Name)
	if err != nil {
		return err
	}

	return nil
}

func (m *SignUpMessage) Size() uint32 {
	return 6 + uint32(m.NameLen)
}

var signUpFormat string = `SignUpMessage{
    Colour:   %v
    NameLen:  %v
    Name:     "%v"
}`

func (m *SignUpMessage) String() string {
	return fmt.Sprintf(signUpFormat, m.PlayerColour, m.NameLen, string(m.Name[:]))
}
