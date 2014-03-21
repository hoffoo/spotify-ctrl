package main

import (
    "flag"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/signal"
    "os/user"
    "strings"
    "syscall"

    clyrics "github.com/hoffoo/go-chartlyrics"
    dbus "github.com/hoffoo/go.dbus"
)

/*
TODO
return an obj instead of a dbus obj
refactor stupid naming
private methods
change dbus to real go.dbus
configurable art_base
*/

var sdbus *dbus.Object
var ART_CACHE string

const ART_BASE = "/.spotify-art/"

func main() {

    var img bool
    flag.BoolVar(&img, "i", false, "get album art image")
    flag.Parse()

    action := flag.Arg(0)

    sdbus = connDbus()
    var song *dbus.Variant
    switch action {
    case "":
        song = Metadata()
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
    case "lyric":
        song = Metadata()

        // buggy spotify dbus only sends a single artist
        songData := song.Value().(map[string]dbus.Variant)
        artist := songData["xesam:artist"].Value().([]string)
        title := songData["xesam:title"]

        Lyric(artist[0], title.String())
    case "next":
        MethodCall("Next")
    case "prev":
        MethodCall("Previous")
    case "pause":
        MethodCall("PlayPause")
    default:
        OpenUri(action)
    }

    // if we supplied the image arg update the album art
    if img {
        u, err := user.Current()
        if err != nil {
            panic(err)
        }

        ART_CACHE = u.HomeDir + ART_BASE
        if song == nil {
            song = Metadata()
        }

        songData := song.Value().(map[string]dbus.Variant)
        artUrl := songData["mpris:artUrl"].Value().(string)
        Art(artUrl)
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

func Art(url string) {

    idx := strings.LastIndex(url, "/")
    filename := url[idx+1:]

    outfile, err := os.OpenFile(ART_CACHE+filename, os.O_RDONLY, 0660)

    if os.IsNotExist(err) {
        sig := make(chan os.Signal, 1)
        signal.Notify(sig, syscall.SIGTERM, syscall.SIGQUIT)

        go func() {
            <-sig
            os.Remove(ART_CACHE + filename)
        }()

        resp, err := http.Get("http://d3rt1990lpmkn.cloudfront.net/unbranded/" + filename)

        if err != nil {
            fmt.Fprintln(os.Stderr, "couldnt download album art")
            return
        }

        outfile, _ = os.Create(ART_CACHE + filename)

        _, ioerr := io.Copy(outfile, resp.Body)
        resp.Body.Close()

        if ioerr != nil {
            outfile.Close()
            os.Remove(ART_CACHE + filename)
            fmt.Fprintln(os.Stderr, "failed getting the album art file")
            return
        }

        outfile.Close()
        close(sig)
    }

    os.Remove(ART_CACHE + "cur")
    os.Link(ART_CACHE+filename, ART_CACHE+"cur")
}

func Lyric(artist, title string) {
    title = strings.Replace(title, "\"", "", -1)
    r, err := clyrics.SearchLyricDirect(artist, title)

    fmt.Printf("Lyrics:  %s, %s\n\n", artist, title)

    if err != nil {
        fmt.Printf("Fetching Failed: %s\n", err)
    } else if strings.Trim(r.Lyric, "") == "" { // TODO better check for no lyrics
        r, err = clyrics.SearchLyric(artist, title, 20)

        if err != nil {
            fmt.Printf("Error getting song update link: %s", err)
        } else {
            fmt.Printf("No lyrics. You should add them here:\n", r.LyricUrl)
            for _, k := range r.SearchLyricResult {
                fmt.Printf("\n%s", k.SongUrl)
            }
        }
    } else {
        fmt.Printf("%s\n", r.Lyric)
    }
}
