package message

import (
	"encoding/binary"
	"fmt"
	"io"
	"slices"
)

var (
	ErrInvalidSig error = fmt.Errorf("invalid signature")
)

type Header struct {
	Signature [4]byte
	MsgType   MessageType
	_         [2]byte
	MsgLen    uint32
}

func NewHeader() *Header {
	return &Header{
		Signature: [4]byte{83, 80, 71, 77}, // SPGM
	}
}

func (h *Header) Read(r io.Reader) error {
	err := binary.Read(r, binary.BigEndian, h)

	if err != nil {
		return err
	}
	if !slices.Equal(h.Signature[:], []byte("SPGM")) {
		return ErrInvalidSig
	}

	return nil
}

func (h *Header) Write(w io.Writer) error {
	if !slices.Equal(h.Signature[:], []byte("SPGM")) {
		return ErrInvalidSig
	}

	return binary.Write(w, binary.BigEndian, h)
}

var headerFormat string = `Header{
    MsgType: %v
    MsgLen: %v
}`

func (h Header) String() string {
	return fmt.Sprintf(headerFormat, h.MsgType, h.MsgLen)
}
