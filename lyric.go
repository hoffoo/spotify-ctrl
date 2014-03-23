package main

import (
	"fmt"
	"strings"

    clyrics "github.com/hoffoo/go-chartlyrics"
)

// Download lyrics from chartlyrics
func Lyric(artist, title string) (err error) {
	title = strings.Replace(title, "\"", "", -1)
	r, err := clyrics.SearchLyricDirect(artist, title)

	fmt.Printf("Lyrics:  %s, %s\n\n", artist, title)

	if err != nil {
		return fmt.Errorf("Fetching lyrics failed: %s\n", err)
	} else if strings.Trim(r.Lyric, "") == "" { // TODO better check for no lyrics
		r, err = clyrics.SearchLyric(artist, title, 20)

		if err != nil {
			return fmt.Errorf("Error getting song update link: %s", err)
		} else {
			return fmt.Errorf("No lyrics. :(")
		}
	} else {
		fmt.Printf("%s\n", r.Lyric)
	}

    return nil
}
