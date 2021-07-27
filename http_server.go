package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func StartHttpServer() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/oauth/spotify", spotifyAuthHandler)
	http.HandleFunc("/nowplaying", nowPlayingRender)
	http.HandleFunc("/nowplaying/ws", nowPlayingWS)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
