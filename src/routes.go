package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
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


    Render("index.html", w, data)
}

func CreateRouter() *mux.Router {
    router := mux.NewRouter()
    router.HandleFunc("/", index).Methods("GET")

    log.Println("Routers created")
    return router
}
