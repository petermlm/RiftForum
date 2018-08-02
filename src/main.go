package main

import (
    "log"
    "net/http"
)

func main() {
    log.Println("Rift Forum starting")

    router := CreateRouter()

    log.Println("Serving")
    status := http.ListenAndServe(":8080", router)
    log.Fatal(status)
}
