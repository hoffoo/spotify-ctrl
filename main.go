package main

import (
	"flag"
	"fmt"
	dbus "github.com/guelfey/go.dbus"
)

func main() {

	flag.Parse()
	action := flag.Arg(0)

	if action == "" {
		fmt.Println("Specify an action")
		return
	}

	switch action {
	case "next":
		Next()
	case "prev":
		Previous()
	case "pause":
		PlayPause()
	case "cur":
		CurSong()
	default:
		fmt.Printf("Invalid action %s\n", action)
		return
	}
}

func connDbus() *dbus.Object {

	conn, err := dbus.SessionBus()

	// couldnt connect to session bus
	if err != nil {
		panic(err)
	}

	return conn.Object("org.mpris.MediaPlayer2.spotify", "/org/mpris/MediaPlayer2")
}

func Next() {
	connDbus().Call("Next", 0)
}

func Previous() {
	connDbus().Call("Previous", 0)
}

func PlayPause() {
	connDbus().Call("PlayPause", 0)
}

func CurSong() {

	sdata   := new(dbus.Variant)	// song data
	pstatus := new(dbus.Variant)	// playing status
	sdbus   := connDbus()

	// get song data, quit on err
	err := sdbus.Call("Get", 0, "org.mpris.MediaPlayer2.Player", "Metadata").Store(sdata)
	if err != nil {
		// most likely spotify not running
		return
	}
	sdbus.Call("Get", 0, "org.mpris.MediaPlayer2.Player", "PlaybackStatus").Store(pstatus)

	// buggy spotify dbus only sends a single artist
	songData := sdata.Value().(map[string]dbus.Variant)
	artist   := songData["xesam:artist"].Value().([]string)
	title    := songData["xesam:title"]
	rating   := int(songData["xesam:autoRating"].Value().(float64) * 100)

	if songStatus := pstatus.Value().(string); songStatus == "Paused" {
		fmt.Printf("(paused) %s %s (paused)", artist[0], title)
	} else {
		fmt.Printf("%s %s (%d)", artist[0], title, rating)
	}
}
