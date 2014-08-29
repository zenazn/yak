package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/nelhage/go.cli/config"
)

const slackURL = "https://%s.slack.com/services/hooks/incoming-webhook?token=%s"

var token, domain, channel, username, icon string

func init() {
	flag.StringVar(&domain, "domain", "", "Slack domain (you.slack.com)")
	flag.StringVar(&token, "token", "", "Secret token for incoming webhook")
	flag.StringVar(&channel, "channel", "#general", "Channel to post in")
	flag.StringVar(&username, "username", "yakbot", "Username to post as")
	flag.StringVar(&icon, "icon", ":speech_balloon:", "Icon URL or emoji")
}

type IncomingWebhook struct {
	Channel   string `json:"channel,omitempty"`
	Username  string `json:"username,omitempty"`
	Text      string `json:"text,omitempty"`
	IconURL   string `json:"icon_url,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
}

func main() {
	config.LoadConfig(flag.CommandLine, "yak.conf")
	flag.Parse()

	if domain == "" {
		fmt.Println("Error: must provide a domain")
		os.Exit(1)
	}
	if token == "" {
		fmt.Println("Error: must provide a token")
		os.Exit(1)
	}

	message := strings.Join(flag.Args(), " ")

	payload := IncomingWebhook{
		Channel:  channel,
		Username: username,
		Text:     message,
	}
	if strings.HasPrefix(icon, ":") {
		payload.IconEmoji = icon
	} else if icon != "" {
		payload.IconURL = icon
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error while marshaling body:", err)
		os.Exit(1)
	}

	u := fmt.Sprintf(slackURL, domain, token)
	response, err := http.PostForm(u, url.Values{
		"payload": []string{string(buf)},
	})

	if err != nil {
		fmt.Println("Error while submitting message:", err)
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		fmt.Println("Warning: received non-200 status",
			response.StatusCode)
	}
	io.Copy(os.Stdout, response.Body)
	fmt.Println()
}
