package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// heavily based on https://github.com/gorilla/websocket/blob/master/examples/echo/server.go

var upgrader = websocket.Upgrader{}

func nowPlayingWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	for {
		nowPlaying := GetNowPlaying()
		err = c.WriteMessage(websocket.TextMessage, []byte(nowPlaying))
		if err != nil {
			panic(err)
		}

		time.Sleep(2 * time.Second)
	}
}

func nowPlayingRender(w http.ResponseWriter, r *http.Request) {
	dat, err := ioutil.ReadFile("now_playing.html")
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(dat))
}
