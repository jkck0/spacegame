package main

import (
	"flag"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/jkck0/spacegame/connection"
	"github.com/jkck0/spacegame/message"
)

func main() {
	debug := flag.Bool("d", false, "enables debug printing")
	flag.Parse()

	u := url.URL{Scheme: "wss", Host: "spacegame.io:443", Path: "/"}
	gc, err := connection.Connect(u, *debug)
	if err != nil {
		panic(err)
	}

	msg := message.NewSignUpMessage()
	msg.PlayerColour = message.PlayerColour{
		R: 232,
		G: 98,
		B: 21,
	}
	msg.Name = []byte("jkck")
	msg.NameLen = uint16(len(msg.Name))

	gc.Send <- msg

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-done:
			break loop
		case <-gc.Recv:
			continue
		}
	}
}
