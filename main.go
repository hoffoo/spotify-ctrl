package main

import (
	dbus "github.com/guelfey/go.dbus"
	introspect "github.com/guelfey/go.dbus/introspect"
	"fmt"
	"flag"
)

func main() {

	flag.Parse()
	action := flag.Arg(0)

	if action == "" {
		fmt.Println("Specify an action")
		return
	}

	switch (action) {
	case "next":
		Next()
	case "prev":
		Previous()
	case "pause":
		PlayPause()
	default:
		fmt.Printf("Invalid action %s\n", action)
		return
	}
}

func spotbus() *dbus.Object {
	conn, err := dbus.SessionBus()

	if err != nil {
		panic(err)
	}

	return conn.Object("org.mpris.MediaPlayer2.spotify", "/org/mpris/MediaPlayer2")
}

func Next() {
	spotbus().Call("Next", 0)
}

func Previous() {
	spotbus().Call("Previous", 0)
}

func PlayPause() {
	spotbus().Call("PlayPause", 0)
}

type metadata interface {
	url() string
}

// FIXME doesnt work
func CurSong() *introspect.Node {
	conn, _ := dbus.SessionBus()

	node, _ := introspect.Call(conn.Object("org.mpris.MediaPlayer2.spotify", "/org/mpris/MediaPlayer2"))

	return node
}
