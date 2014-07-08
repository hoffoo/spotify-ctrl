package main

import (
	"github.com/godbus/dbus"
	"strings"
)

var sdbus *dbus.Object

type Spotify struct {
	Artist  string
	Title   string
	Rating  int
	Status  string
	Url     string
	ArtUrl  string
	ArtFile string
}

// Updates the spotify struct with current data from spotify
// TODO return err here
func (c *Spotify) Get() {
	if sdbus == nil {
		sdbus = connDbus()
	}

	song := Metadata()
	pstatus := Status()

	// TODO buggy spotify only sends a single artist
	songData := song.Value().(map[string]dbus.Variant)
	c.Artist = songData["xesam:artist"].Value().([]string)[0]
	c.Title = songData["xesam:title"].Value().(string)
	c.Rating = int(songData["xesam:autoRating"].Value().(float64) * 100)
	c.Status = pstatus.Value().(string)
	c.Url = songData["xesam:url"].Value().(string)
	c.ArtUrl = songData["mpris:artUrl"].Value().(string)

	idx := strings.LastIndex(c.ArtUrl, "/")
	c.ArtFile = c.ArtUrl[idx+1:]
}

func SpotifyMethod(method string) error {
	if sdbus == nil {
		sdbus = connDbus()
	}

	c := sdbus.Call(method, 0)
	return c.Err
}

func OpenUri(uri string) {
	sdbus.Call("OpenUri", 0, uri)
}

// TODO return err here
func Metadata() *dbus.Variant {

	// get song data, quit on err
	song, err := sdbus.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")
	if err != nil {
		panic(err) // most likely spotify not running
	}

	return &song
}

func Status() *dbus.Variant {

	pstatus, err := sdbus.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")

	if err != nil {
		panic(err) // most likely spotify not running
	}

	return &pstatus
}

func connDbus() *dbus.Object {

	conn, err := dbus.SessionBus()

	// couldnt connect to session bus
	if err != nil {
		panic(err)
	}

	return conn.Object("org.mpris.MediaPlayer2.spotify", "/org/mpris/MediaPlayer2")
}
