package main

import (
	"fmt"
	"log"
)

func RiftForumPanic(msg string, err error) {
	log.Fatal(err)
	panic(msg)
}

func MakeBaseURL() string {
	var protocol string

	if Config.HTTPS {
		protocol = "https"
	} else {
		protocol = "http"
	}

	return fmt.Sprintf("%s://%s", protocol, Config.BaseURL)
}

func RiftLink(url string) string {
	return fmt.Sprintf("%s%s", Config.APIBase, url)
}

func MakeRedditLink(subreddit string) string {
	return fmt.Sprintf("https://reddit.com%s", subreddit)
}
