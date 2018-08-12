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

func topics(writer http.ResponseWriter, r *http.Request) {
    var err error
    db := GetDBCon()
    form_title := r.PostFormValue("title")
    form_message := r.PostFormValue("message")

    user := new(User)
    err = db.Model(user).Where("Username = ?", "admin").Select()

    if err != nil {
        panic(err)
    }

    // Topic
    topic := &Topic{
        Title: form_title,
        Author: user,
        AuthorId: user.Id,
    }

    err = db.Insert(topic)

    if err != nil {
        panic(err)
    }

    // Message
    message := &Message{
        Message: form_message,
        Author: user,
        AuthorId: user.Id,
        Topic: topic,
        TopicId: topic.Id,
    }

    err = db.Insert(message)

    if err != nil {
        panic(err)
    }

    http.Redirect(writer, r, "/", http.StatusSeeOther)
}

func CreateRouter() *mux.Router {
    router := mux.NewRouter()
    router.HandleFunc("/", index).Methods("GET")
    router.HandleFunc("/topics", topics).Methods("POST")

    router.
        PathPrefix("/static/style/").
        Handler(http.StripPrefix("/static/style/", http.FileServer(http.Dir("../static/style/"))))

    log.Println("Routers created")
    return router
}
