package main

import "log"

func RiftForumPanic(msg string, err error) {
    log.Fatal(err)
    panic(msg)
}
