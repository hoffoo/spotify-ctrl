package main

import (
	"flag"
	"fmt"
	dbus "github.com/hoffoo/go.dbus"
	"io"
	"net/http"
	"os"
	"os/user"
	"strings"
)

var sdbus *dbus.Object

const ART_BASE = "/.spotify-art/"
var ART_CACHE string

func main() {

	var img bool
	flag.BoolVar(&img, "i", false, "get album art image")
	flag.Parse()

	action := flag.Arg(0)

	if img {
		u, err := user.Current()
		if err != nil{
			panic(err)
		}
		ART_CACHE = u.HomeDir+ART_BASE
	}

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

		if img == true {
			artUrl := songData["mpris:artUrl"].Value().(string)
			Art(artUrl)
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

func Art(url string) {

	idx := strings.LastIndex(url, "/")
	filename := url[idx+1:]

	outfile, err := os.OpenFile(ART_CACHE+filename, os.O_RDONLY, 0660)

	if os.IsNotExist(err) {
		resp, err := http.Get("https://d3rt1990lpmkn.cloudfront.net/unbranded/" + filename)
		defer resp.Body.Close()

		if err != nil {
			fmt.Fprintln(os.Stderr, "couldnt download album art")
			return
		}

		outfile, err = os.Create(ART_CACHE+filename)

		io.Copy(outfile, resp.Body)
	}

	os.Remove(ART_CACHE+"cur")
	os.Link(ART_CACHE+filename, ART_CACHE+"cur")
}
