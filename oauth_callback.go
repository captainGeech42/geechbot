package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func spotifyAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body style=\"font-size: 16px;\"><pre>Please return back to geechbot</pre></body></html>")

	go HandleSpotifyOauthCallback(r)
}

func StartHttpServer() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/oauth/spotify", spotifyAuthHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
