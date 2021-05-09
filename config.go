package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Channel    string `yaml:"channel"`
	Username   string `yaml:"username"`
	OAuthToken string `yaml:"oauth_token"`
}

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
