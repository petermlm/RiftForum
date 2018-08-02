package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
    data := struct {
        Username string
        Usertype string
        Items []string
    }{
        Username: "Peter",
        Usertype: "Admin",
        Items: []string{
            "One item",
            "Another item",
            "Final item",
        },
    }

    Render("index.html", w, data)
}

func CreateRouter() *mux.Router {
    router := mux.NewRouter()
    router.HandleFunc("/", index).Methods("GET")

    log.Println("Routers created")
    return router
}
