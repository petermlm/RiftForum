package main

import (
    "log"
    "net/http"
)

func main() {
    log.Println("Rift Forum Starting")

    InitDB()
    defer CloseDB()
    InitTmpl()
    InitAuth()

    log.Println("Starting Server")
    router := CreateRouter()
    status := http.ListenAndServe(HostAndPort, router)
    log.Fatal(status)
}
