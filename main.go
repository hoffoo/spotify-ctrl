package main

import (
	dbus "github.com/guelfey/go.dbus"
	dbusprop "github.com/guelfey/go.dbus/prop"
	"fmt"
	"flag"
	"regexp"
)

func main() {

	flag.Parse()
	action := flag.Arg(0)

	if action == "" {
		fmt.Println("Specify an action")
		return
	}

	if ok, _ := regexp.MatchString("next|prev|pause", action); ok == false {
		fmt.Printf("Invalid action %s\n", action)
		return
	}

	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	sdbus := conn.Object("org.mpris.MediaPlayer2.spotify", "/org/mpris/MediaPlayer2")

	switch (action) {
	case "next":
		Next(sdbus)
	case "prev":
		Previous(sdbus)
	case "pause":
		PlayPause(sdbus)
	}
}

func Next(o *dbus.Object) *dbus.Call {
	return o.Call("Next", 0)
}

func Previous(o *dbus.Object) *dbus.Call {
	return o.Call("Previous", 0)
}

func PlayPause(o *dbus.Object) *dbus.Call {
	return o.Call("PlayPause", 0)
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
