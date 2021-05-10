package main

import (
	"fmt"

	"github.com/gempir/go-twitch-irc/v2"
)

func main() {
	config = LoadConfig()

	client := twitch.NewClient(config.Twitch.Username, config.Twitch.OAuthToken)

	client.OnConnect(func() {
		fmt.Println("Successfully connected to TMI service")
	})

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Printf("Got a message from %s: %s\n", message.User.DisplayName, message.Message)

		if message.Message == "!ping" {
			client.Say(config.Twitch.Channel, "pong")
		}
	})

	client.Join(config.Twitch.Channel)

	fmt.Println("Starting HTTP server")
	go StartHttpServer()
	fmt.Println("Sucessfully started HTTP server")

	fmt.Println("Authing Spotify")
	AuthToSpotify()

	fmt.Println("Connecting to TMI service")
	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
