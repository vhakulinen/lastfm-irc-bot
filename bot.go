package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/thoj/go-ircevent"
)

const (
	// FeedURLStart is beginning of the url which needs username and
	// FeedURLEnd added
	FeedURLStart = "http://ws.audioscrobbler.com/1.0/user/"
	// FeedURLEnd is meant to be used with FeedURLEdn
	FeedURLEnd = "/recenttracks.rss"
)

// Item is object in the result
type Item struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
}

// Channel is object in the result
type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Items   []Item   `xml:"item"`
}

// Rss is the result object
type Rss struct {
	XMLName  xml.Name  `xml:"rss"`
	Channels []Channel `xml:"channel"`
}

func lastFMCallback(event *irc.Event) {
	username := event.Nick
	msg := strings.Split(event.Message(), " ")
	if msg[0] != ",np" {
		return
	}
	if len(msg) == 2 {
		username = msg[1]
	}
	resp, err := http.Get(FeedURLStart + username + FeedURLEnd)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	print(resp.StatusCode)
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	v := &Rss{}
	err = xml.Unmarshal(content, v)
	if err != nil {
		return
	}
	if len(v.Channels) > 0 && len(v.Channels[0].Items) > 0 {
		event.Connection.Privmsg("#nixers",
			fmt.Sprintf("%s is listening %s", username, v.Channels[0].Items[0].Title))
	}
}

func main() {
	icon := irc.IRC("gnp", "gnp")
	//icon.UseTLS = true
	icon.Debug = true
	err := icon.Connect("irc.nixers.net:6667")
	if err != nil {
		log.Fatal(err)
	}
	icon.AddCallback("001", func(e *irc.Event) {
		icon.Join("#nixers")
	})
	icon.AddCallback("PRIVMSG", lastFMCallback)
	icon.Loop()
}
