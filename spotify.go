package main

import (
	"fmt"
	"net/http"

	"github.com/zmb3/spotify"
)

var spotifyAuthRedirectUri = "http://localhost:8080/oauth/spotify"

var spotifyAuthClient spotify.Authenticator
var spotifyClient spotify.Client

type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func InitSpotifyAuth() {
	spotifyAuthClient = spotify.NewAuthenticator(spotifyAuthRedirectUri, spotify.ScopeUserReadCurrentlyPlaying)
	spotifyAuthClient.SetAuthInfo(config.Spotify.ClientID, config.Spotify.ClientSecret)
	url := spotifyAuthClient.AuthURL("asdf")

	fmt.Printf("Please go to the following URL in your web browser: %s\n", url)
}

func HandleSpotifyOauthCallback(r *http.Request) {
	token, err := spotifyAuthClient.Token("asdf", r)
	if err != nil {
		panic(err)
	}

	spotifyClient = spotifyAuthClient.NewClient(token)

	fmt.Println("Successfully authed to Spotify")
}

func GetNowPlaying() string {
	currentlyPlaying, err := spotifyClient.PlayerCurrentlyPlaying()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Now Playing: %s by %s", currentlyPlaying.Item.Name, currentlyPlaying.Item.Artists[0].Name)
}
