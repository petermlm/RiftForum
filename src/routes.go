package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func index(writer http.ResponseWriter, r *http.Request) {
    topics := GetTopics()
    ser_topics := SerializeTopics(topics)

    data := struct {
        Topics []interface{}
    }{
        Topics: ser_topics,
    }

    Render(writer, "index.html", data)
}

func CreateRouter() *mux.Router {
    router := mux.NewRouter()
    router.HandleFunc("/", index).Methods("GET")

    router.
        PathPrefix("/static/style/").
        Handler(http.StripPrefix("/static/style/", http.FileServer(http.Dir("../static/style/"))))

    log.Println("Routers created")
    return router
}
