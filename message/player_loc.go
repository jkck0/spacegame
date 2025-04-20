package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Player struct {
	ID   uint32
	Name []byte
	PlayerColour
	SecX  uint16
	SecY  uint16
	MapX  float64
	MapY  float64
	MapVX float64
	MapVY float64
}

var playerFormat string = `Player{
    ID:            %v
    Name:          "%v" %v
    PlayerColour:  %v
    SecX:          %v
    SecY:          %v
    MapX:          %v
    MapY:          %v
    MapVX:         %v
    MapVY:         %v
}`

func (p Player) String() string {
	return fmt.Sprintf(
		playerFormat,
		p.ID,
		string(p.Name),
		p.Name,
		p.PlayerColour,
		p.SecX,
		p.SecY,
		p.MapX,
		p.MapY,
		p.MapVX,
		p.MapVY,
	)
}

func (p *Player) Read(r io.Reader) error {
	err := binary.Read(r, binary.BigEndian, &p.ID)
	if err != nil {
		return err
	}

	p.Name = make([]byte, 17)
	_, err = io.ReadFull(r, p.Name)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.BigEndian, &p.PlayerColour)
	if err != nil {
		return err
	}

	// the documentation lies about the endianness of these
	err = binary.Read(r, binary.LittleEndian, &p.SecX)
	if err != nil {
		return err
	}
	io.CopyN(io.Discard, r, 2)

	err = binary.Read(r, binary.LittleEndian, &p.SecY)
	if err != nil {
		return err
	}
	io.CopyN(io.Discard, r, 2)

	err = binary.Read(r, binary.LittleEndian, &p.MapX)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.LittleEndian, &p.MapY)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.BigEndian, &p.MapVX)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.BigEndian, &p.MapVY)
	if err != nil {
		return err
	}

	return nil
}

func (p Player) Size() uint32 {
	return 60
}

type PlayerLocMessage struct {
	NumPlayers uint16
	Players    []Player
}

func NewPlayerLocMessage() *PlayerLocMessage {
	return &PlayerLocMessage{}
}

func (m *PlayerLocMessage) Read(r io.Reader, read uint32) error {
	err := binary.Read(r, binary.BigEndian, &m.NumPlayers)
	if err != nil {
		return err
	}

	m.Players = make([]Player, m.NumPlayers)
	for i := range m.Players {
		err := (&m.Players[i]).Read(r)
		if err != nil {
			return err
		}
	}

	rest := make([]byte, read-m.Size())
	_, err = io.ReadFull(r, rest)
	if err != nil {
		return err
	}

	rest = rest[2:] // dunno why i have to do this
	for i := range m.Players {
		// if the name is not null terminated, there is extra to get at the end
		if m.Players[i].Name[16] != 0 {
			name_rest := rest[:findNull(rest)]
			m.Players[i].Name = append(m.Players[i].Name, name_rest...)
			rest = rest[findNull(rest)+1:]
		}
	}

	return nil
}

// returns the index of the first null byte in the slice, -1 if none
func findNull(buf []byte) int {
	for i, b := range buf {
		if b == 0 {
			return i
		}
	}

	return -1
}

func (m *PlayerLocMessage) Write(w io.Writer) error {
	return nil
}

func (m *PlayerLocMessage) Size() uint32 {
	return 2 + Player{}.Size()*uint32(m.NumPlayers)
}

var playerLocFormat string = `PlayerLocMessage{
    NumPlayers:  %v
    Players:     %v
}`

func (m *PlayerLocMessage) String() string {
	return fmt.Sprintf(playerLocFormat, m.NumPlayers, m.Players)
}
