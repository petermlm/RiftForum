package main

import (
    "log"
    "net/http"
)

func main() {
    log.Println("Rift Forum starting")

    InitDB()
    defer CloseDB()
    InitTmpl()
    router := CreateRouter()

    log.Println("Serving")
    status := http.ListenAndServe(":8080", router)
    log.Fatal(status)
}
