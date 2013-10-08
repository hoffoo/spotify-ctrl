package main

import (
	"flag"
	"fmt"
	dbus "github.com/guelfey/go.dbus"
	"strings"
)

func main() {

	flag.Parse()
	action := flag.Arg(0)

	switch action {
	case "":
		song, pstatus, err := CurSong()

		if err != "" {
			fmt.Printf(err)
			return
		}

		// buggy spotify dbus only sends a single artist
		songData := song.Value().(map[string]dbus.Variant)
		artist := songData["xesam:artist"].Value().([]string)
		title := songData["xesam:title"]
		rating := int(songData["xesam:autoRating"].Value().(float64) * 100)

		if songStatus := pstatus.Value().(string); songStatus == "Paused" {
			fmt.Printf("(paused) %s %s (paused)", artist[0], title)
		} else {
			fmt.Printf("%s %s (%d)", artist[0], title, rating)
		}
	case "url":
		song, _, err := CurSong()

		if err != "" {
			fmt.Printf(err)
			return
		}

		songData := song.Value().(map[string]dbus.Variant)
		url := songData["xesam:url"].Value().(string)
		fmt.Printf("http://open.spotify.com/track/%s", strings.Split(url, ":")[2])
	case "next":
		connDbus().Call("Next", 0)
	case "prev":
		connDbus().Call("Previous", 0)
	case "pause":
		connDbus().Call("PlayPause", 0)
	default:
		connDbus().Call("OpenUri", 0, action)
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

func CurSong() (*dbus.Variant, *dbus.Variant, string) {

	song := new(dbus.Variant)    // song data
	pstatus := new(dbus.Variant) // playing status
	sdbus := connDbus()

	// get song data, quit on err
	err := sdbus.Call("Get", 0, "org.mpris.MediaPlayer2.Player", "Metadata").Store(song)
	if err != nil {
		// most likely spotify not running
		return nil, nil, "Couldnt send to Dbus - is spotify running?"
	}
	sdbus.Call("Get", 0, "org.mpris.MediaPlayer2.Player", "PlaybackStatus").Store(pstatus)

	return song, pstatus, ""
}
