package main

import (
	"io/ioutil"

	"github.com/forestgiant/sliceutil"
	"github.com/gempir/go-twitch-irc/v2"
	"gopkg.in/yaml.v2"
)

type Response struct {
	Commands []string `yaml:"commands"`
	Messages []string `yaml:"messages"`
}

type Responses struct {
	StaticResponses []Response `yaml:"static_responses"`
}

var staticResponses *Responses

func LoadStaticResponses() {
	staticResponses = &Responses{}

	raw, err := ioutil.ReadFile("responses.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(raw, &staticResponses)
	if err != nil {
		panic(err)
	}

	// make sure all messages are under 500 char
	for _, r := range staticResponses.StaticResponses {
		for _, m := range r.Messages {
			if len(m) > 500 {
				panic("Message is too long!: " + m)
			}
		}
	}
}

func HandleStaticResponse(msg twitch.PrivateMessage) *[]string {
	// check if the responses have been loaded into memory
	if staticResponses == nil {
		LoadStaticResponses()
	}

	// find a response for the command
	for _, r := range staticResponses.StaticResponses {
		if sliceutil.Contains(r.Commands, msg.Message) {
			// found a command, send the response
			return &r.Messages
		}
	}

	// no command found
	return nil
}
