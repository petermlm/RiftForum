package main

import (
    "fmt"
    "log"
)

func RiftForumPanic(msg string, err error) {
    log.Fatal(err)
    panic(msg)
}

func MakeBaseUrl() string {
    var protocol string

    if Https {
        protocol = "https"
    } else {
        protocol = "http"
    }

    return fmt.Sprintf("%s://%s", protocol, BaseUrl)
}

func MakeRedditLink(subreddit string) string {
    return fmt.Sprintf("https://reddit.com%s", subreddit)
}
