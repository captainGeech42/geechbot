package main

import (
	"fmt"

	"github.com/gempir/go-twitch-irc/v2"
)

func main() {
	config := LoadConfig()

	client := twitch.NewClient(config.Username, config.OAuthToken)

	client.OnConnect(func() {
		fmt.Println("Connected to TMI service")
	})

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Printf("Got a message from %s: %s\n", message.User.DisplayName, message.Message)

		if message.Message == "!ping" {
			client.Say(config.Channel, "pong")
		}
	})

	client.Join(config.Channel)

	fmt.Println("Authenticating to TMI service")
	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
