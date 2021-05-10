package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var spotifyAuthRedirectUri = "http://localhost:8080/oauth/spotify"

type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func AuthToSpotify() {
	/*
		scope: user-read-currently-playing
	*/

	url := fmt.Sprintf("https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&state=%s&scope=%s&show_dialog=%s",
		config.Spotify.ClientID,                 // client id
		url.QueryEscape(spotifyAuthRedirectUri), // redirect uri
		"asdf",                                  // state
		"user-read-currently-playing",           // scope
		"false",                                 // show dialog
	)

	fmt.Printf("Please go to the following URL in your web browser: %s\n", url)
}

func HandleSpotifyOauthCallback(r *http.Request) {
	u, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		panic(err)
	}

	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		panic(err)
	}

	code := m["code"]

	if len(code) == 0 {
		// didn't get a code, probably error
		err_msg := m["error"]
		if len(err_msg) == 0 {
			fmt.Printf("Failed to parse Spotify callback: %s\n", r.RequestURI)
			os.Exit(1)
		}

		fmt.Printf("Got an error from Spotify auth callback: %s", err_msg[0])
		os.Exit(1)
	}

	fmt.Printf("Spotify auth code: %s\n", code[0])

	// request the refresh/access tokens
	postBody := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code[0]},
		"redirect_uri": {spotifyAuthRedirectUri},
	}

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(postBody.Encode()))
	if err != nil {
		panic(err)
	}

	fmt.Println("Basic " + base64.StdEncoding.EncodeToString([]byte(config.Spotify.ClientID+":"+config.Spotify.ClientSecret)))

	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(config.Spotify.ClientID+":"+config.Spotify.ClientSecret)))

	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Got a bad status code in token request: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	respFields := SpotifyTokenResponse{}

	fmt.Println(string(respBytes))

	err = json.Unmarshal(respBytes, &respFields)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Spotify access token: %s\n", respFields.AccessToken)
}
