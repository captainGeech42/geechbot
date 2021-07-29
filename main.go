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
		} else if message.Message == "!np" {
			client.Say(config.Twitch.Channel, GetNowPlaying())
		} else if message.Message == "!clear" {
			_, ok := message.User.Badges["moderator"]
			if !ok {
				_, ok = message.User.Badges["broadcaster"]
			}

			if ok {
				client.Say(config.Twitch.Channel, "/clear")
			}
		} else if message.Message == "!reload_commands" {
			_, ok := message.User.Badges["moderator"]
			if !ok {
				_, ok = message.User.Badges["broadcaster"]
			}

			if ok {
				LoadStaticResponses()
				client.Say(config.Twitch.Channel, "commands reloaded!")
			}
		} else {
			// handle as a static response from the yaml file
			resp := HandleStaticResponse(message)
			if resp != nil {
				for _, m := range *resp {
					client.Say(config.Twitch.Channel, m)
				}
			}
		}
	})

	client.Join(config.Twitch.Channel)

	fmt.Println("Starting HTTP server")
	go StartHttpServer()
	fmt.Println("Sucessfully started HTTP server")

	fmt.Println("Authing Spotify")
	InitSpotifyAuth()

	fmt.Println("Connecting to TMI service")
	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
