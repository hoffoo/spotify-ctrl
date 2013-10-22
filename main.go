package main

import (
	//dbus "../go.dbus"
	dbus "github.com/hoffoo/go.dbus"
	"flag"
	"fmt"
	//"os"
	"strings"
)

var sdbus *dbus.Object

const ART_CACHE = "~/.spotify-art/"

func main() {

	var img bool
	flag.BoolVar(&img, "i", false, "get album art image")
	flag.Parse()

	action := flag.Arg(0)

	sdbus = connDbus()
	switch action {
	case "":
		song := Metadata()
		pstatus := Status()

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
		song := Metadata()

		songData := song.Value().(map[string]dbus.Variant)
		url := songData["xesam:url"].Value().(string)
		fmt.Printf("http://open.spotify.com/track/%s\n", strings.Split(url, ":")[2])
	case "next":
		MethodCall("Next")
	case "prev":
		MethodCall("Previous")
	case "pause":
		MethodCall("PlayPause")
	default:
		OpenUri(action)
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

func MethodCall(method string) {
	sdbus.Call(method, 0)
}

func OpenUri(uri string) {
	sdbus.Call("OpenUri", 0, uri)
}

func Metadata() *dbus.Variant {

	// get song data, quit on err
	song, err := sdbus.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")
	if err != nil {
		panic(err) // most likely dbus not running
	}

	return &song
}

func Status() *dbus.Variant {

	pstatus, err := sdbus.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")

	if err != nil {
		panic(err) // most likely dbus not running
	}

	return &pstatus
}

//func FetchArt(url string) {
//
//	idx := strings.LastIndex(url, "/")
//	filename := url[idx+1:]
//
//}
