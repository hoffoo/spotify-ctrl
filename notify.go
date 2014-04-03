package main

import (
    "io/ioutil"
    "os"
    "time"

    "github.com/mqu/go-notify"
)

func Notify(s *Spotify) {
    f, err := os.OpenFile(ART_CACHE+"track", os.O_RDWR|os.O_CREATE, 0600)

    if err != nil {
        panic(err)
    }

    lstTrackB, err := ioutil.ReadAll(f)

    if err != nil {
        panic(err)
    }

    lstTrack := string(lstTrackB)
    songId := s.Artist + " - " + s.Title

    // song hasnt changed
    if lstTrack == songId {
        return
    }

    f.Truncate(0)
    f.WriteAt([]byte(songId),0)
    f.Close()
    doNotify(s)
}

func doNotify(s *Spotify) {
    notify.Init("Spotify")
    nowPlaying := notify.NotificationNew(s.Artist+" - "+s.Title, "Now Playing", ART_CACHE+s.ArtFile)

    if nowPlaying == nil {
        panic("now -playing null")
    }

    notify.NotificationSetTimeout(nowPlaying, 3000)
    notify.NotificationShow(nowPlaying)

    time.Sleep(3000 * 1000000)
    notify.NotificationClose(nowPlaying)
    notify.UnInit()
}
