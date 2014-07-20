spotify-ctrl
==========

A script to control Spotify over dbus - good for keybinding and getting album
artwork.

```sh
# print current song: 'artist "track" (spotify rating)' - I use it for conky
spotify-ctrl [-i]

spotify-ctrl next
spotify-ctrl prev
spotify-ctrl pause		# pause toggle
spotify-ctrl url 			# print the current song url

# open a url
spotify-ctrl http://open.spotify.com/track/1ipS1pdAnpqTz0QMZePTz1

```

Passing -i will also attempt to download a high rez album art image from spotify

Images are stored in ~/.spotify-art/ , make this directory ahead of time.  If
the image does not exist it is downloaded.  If it exists a hard link 'cur' is 
made to the appropriate image. 

My conky lines looks like this: 
```conky
${execi 10 /opt/bin/go-spotify -i}
${image ~/.spotify-art/cur -n -p 0,27 -s 286x286}
```

