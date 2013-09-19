package main

import (
	dbus "github.com/guelfey/go.dbus"
	dbusprop "github.com/guelfey/go.dbus/prop"
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

// FIXME doesnt work
func CurSong(c *dbus.Conn) (dbus.Variant, error) {
	data := map[string]map[string]*dbusprop.Prop {
		"org.mpris.MediaPlayer2.spotify": {
			"Metadata": {
			},
		},
	}
	prop := dbusprop.New(c, "org.freedesktop.DBus.Properties.Get", data)

	return prop.Get("org.mpris.MediaPlayer2.Player", "Metadata")
}
