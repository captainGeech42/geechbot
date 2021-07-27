package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var spotifyAuthRedirectUri = "http://localhost:8080/oauth/spotify"

var spotifyAuthClient spotify.Authenticator
var spotifyClient spotify.Client

const cacheFilePath = ".spotify_auth_cache"

type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func cacheAuth(token *oauth2.Token) {
	b, err := json.Marshal(token)
	if err != nil {
		panic(err)
	}

	if err = os.WriteFile(cacheFilePath, b, 0600); err != nil {
		panic(err)
	}
}

func readAuthFromCache() *oauth2.Token {
	// check if there is a cache file
	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		// no cache
		return nil
	}

	// read in cached token
	b, err := ioutil.ReadFile(cacheFilePath)
	if err != nil {
		panic(err)
	}

	// unmarshal it
	var token oauth2.Token
	if err = json.Unmarshal(b, &token); err != nil {
		panic(err)
	}

	// make sure the token is still valid
	if token.Expiry.After(time.Now()) {
		return nil
	}

	// token valid, good to go
	return &token
}

func InitSpotifyAuth() {
	// try to read a cached token
	token := readAuthFromCache()
	if token == nil {
		spotifyClient = spotifyAuthClient.NewClient(token)
		fmt.Println("Successfully authed to Spotify using cached creds")

		return
	}

	// no cache, auth via normal oauth flow
	// creds will be cached

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

	// cache token
	cacheAuth(token)

	spotifyClient = spotifyAuthClient.NewClient(token)

	fmt.Println("Successfully authed to Spotify")
}

func GetNowPlaying() string {
	// TODO: need to verify that we are authed and spotifyClient isn't null

	currentlyPlaying, err := spotifyClient.PlayerCurrentlyPlaying()
	if err != nil {
		panic(err)
	}

	if currentlyPlaying == nil {
		return "No song playing"
	}

	return fmt.Sprintf("Now Playing: %s by %s", currentlyPlaying.Item.Name, currentlyPlaying.Item.Artists[0].Name)
}

func spotifyAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body style=\"font-size: 16px;\"><pre>Please return back to geechbot</pre></body></html>")

	go HandleSpotifyOauthCallback(r)
}
