package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func index(writer http.ResponseWriter, r *http.Request) {
    user := new(User)
    db := GetDBCon()
    err := db.Model(user).
        Where("username = ?", "admin").
        Select()

    if err != nil {
        panic(err)
    }

    data := struct {
        Username string
        Usertype string
        Items []string
    }{
        Username: user.Username,
        Usertype: user.GetUserType(),
        Items: []string{
            "One item",
            "Another item",
            "Final item",
        },
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
