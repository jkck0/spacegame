package connection

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/jkck0/spacegame/message"
)

type GameConn struct {
	conn *websocket.Conn
	Recv <-chan message.Message
	Send chan<- message.Message
	// used for the readHandler to tell the writeHandler to pong
	ping chan int
}

func Connect(u url.URL, debug bool) (*GameConn, error) {
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return &GameConn{}, err
	}
	log.Print("dial: connection created")

	recv := make(chan message.Message)
	send := make(chan message.Message)
	gc := &GameConn{
		conn: c,
		Recv: recv,
		Send: send,
	}

	go gc.readHandler(recv, debug)
	go gc.writeHandler(send)

	return gc, nil
}

func (gc *GameConn) readHandler(out chan<- message.Message, debug bool) {
	for {
		_, r, err := gc.conn.NextReader()
		if err != nil {
			fmt.Println(err)
			return
		}

		if debug {
			raw, err := io.ReadAll(r)
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println(hex.Dump(raw))
			r = bytes.NewReader(raw)
		}

		msg, err := message.ReadMessage(r)
		if err != nil {
			fmt.Println(err)
			continue
		}

		log.Printf("recieve: %v", msg)
		if _, ok := msg.(*message.PingMessage); ok {
			gc.ping <- 1
			continue
		}

		out <- msg
	}
}

func (gc *GameConn) writeHandler(in <-chan message.Message) {
	pong := message.NewPongMessage()

	for {
		select {
		case <-gc.ping:
			err := gc.sendMessage(pong)
			if err != nil {
				fmt.Println(err)
			}

		case msg := <-in:
			err := gc.sendMessage(msg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (gc *GameConn) sendMessage(msg message.Message) error {
	log.Printf("send: %v", msg)

	w, err := gc.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}

	err = message.WriteMessage(w, msg)
	if err != nil {
		return err
	}
	w.Close()

	return nil
}

func (gc *GameConn) Close() {
	gc.conn.Close()
}
