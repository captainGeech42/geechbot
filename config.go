package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Twitch struct {
	Channel    string `yaml:"channel"`
	Username   string `yaml:"username"`
	OAuthToken string `yaml:"oauth_token"`
}

type Spotify struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type Config struct {
	Twitch  Twitch  `yaml:"twitch"`
	Spotify Spotify `yaml:"spotify"`
}

var config Config

func LoadConfig() Config {
	config := Config{}

	raw, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(raw, &config)
	if err != nil {
		panic(err)
	}

	return config
}
