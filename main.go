package main

import (
    "flag"
    "fmt"
    "os/user"
    "strings"
    "os"
)

/*
TODO
change dbus to real go.dbus
configurable art_base?
*/
var ART_CACHE string
const ART_BASE = "/.spotify-art/"

func main() {

    var img bool
    flag.BoolVar(&img, "i", false, "get album art image")
    flag.Parse()

    action := flag.Arg(0)
    switch action {
    case "next":
        SpotifyMethod("Next")
    case "prev":
        SpotifyMethod("Previous")
    case "pause":
        SpotifyMethod("PlayPause")
    }

    S := new(Spotify)
    S.Get()

    switch action {
    case "":
        if S.Status == "Paused" {
            fmt.Printf("(paused) %s - \"%s\" (paused)", S.Artist, S.Title)
        } else {
            fmt.Printf("%s - \"%s\" (%d)", S.Artist, S.Title, S.Rating)
        }
    case "url":
        fmt.Printf("http://open.spotify.com/track/%s\n", strings.Split(S.Url, ":")[2])
    case "lyric":
        Lyric(S.Artist, S.Title)
    }

    // if we supplied the -i arg update the album art
    u, err := user.Current()
    if err != nil {
        panic(err)
    }

    ART_CACHE = u.HomeDir + ART_BASE
    if img {
        err = GetArt(ART_CACHE, S.ArtUrl)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    }

}

