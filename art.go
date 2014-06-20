package main

import (
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Fetches album art or links a file named "cur" to the current art.
// If this filename (taken from the url) has not been downloaded
// we download from http://d3rt1990lpmkn.cloudfront.net/unbranded/
// Once we have the image we hard link to "cur". If we already have
// the image filename its only hard linked without downloading.
//
// If something goes wrong while downloading (SIGTERM, SIGQUIT)
// it removes the file so we dont get partial imges.
func GetArt(artCache, url string) (err error) {

	idx := strings.LastIndex(url, "/")
	filename := url[idx+1:]

	outfile, err := os.OpenFile(artCache+filename, os.O_RDONLY, 0660)

	if os.IsNotExist(err) {

		// catch SIGTERM and cleanup temporary file (may not have read
		// the whole image yet)
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			<-sig
			os.Remove(artCache + filename)
		}()

		resp, err := http.Get("http://d3rt1990lpmkn.cloudfront.net/unbranded/" + filename)

		if err != nil {
			return err
		}

		outfile, _ = os.Create(artCache + filename)

		_, ioerr := io.Copy(outfile, resp.Body)
		resp.Body.Close()

		if ioerr != nil {
			outfile.Close()
			os.Remove(artCache + filename)
			return ioerr
		}

		outfile.Close()
		close(sig)
	}

	// remove the old hard link
	os.Remove(artCache + "cur")
	// link to the new current album art
	os.Link(artCache+filename, artCache+"cur")

	return nil
}
