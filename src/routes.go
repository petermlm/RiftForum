package main

import (
    "fmt"
    "strconv"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func index(writer http.ResponseWriter, r *http.Request) {
    topics := GetTopics()
    data := SerializeTopics(topics)
    Render(writer, "index.html", data)
}

func topics_post(writer http.ResponseWriter, r *http.Request) {
    var err error
    db := GetDBCon()
    form_title := r.PostFormValue("title")
    form_message := r.PostFormValue("message")

    // Author
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

func topic_get(writer http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    topic_id_parsed, err := strconv.ParseUint(vars["id"], 10, 32)

    if err != nil {
        // TODO
    }

    topic_id := uint(topic_id_parsed)
    topic := GetTopic(topic_id)
    data := SerializeTopic(topic)
    Render(writer, "topic.html", data)
}

func topic_post(writer http.ResponseWriter, r *http.Request) {
    var err error
    db := GetDBCon()
    form_message := r.PostFormValue("message")

    // Author
    user := new(User)
    err = db.Model(user).Where("Username = ?", "admin").Select()

    if err != nil {
        panic(err)
    }

    // Topic
    vars := mux.Vars(r)
    topic_id_parsed, err := strconv.ParseUint(vars["id"], 10, 32)

    if err != nil {
        // TODO
    }

    topic_id := uint(topic_id_parsed)
    topic := GetTopic(topic_id)
    UpdateTopic(topic)

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

    redirect_path := fmt.Sprintf("/topics/%d", topic_id)
    http.Redirect(writer, r, redirect_path, http.StatusSeeOther)
}

func CreateRouter() *mux.Router {
    router := mux.NewRouter()
    router.HandleFunc("/", index).Methods("GET")
    router.HandleFunc("/topics", topics_post).Methods("POST")
    router.HandleFunc("/topics/{id:[0-9]+}", topic_get).Methods("GET")
    router.HandleFunc("/topics/{id:[0-9]+}", topic_post).Methods("POST")

    router.
        PathPrefix("/static/style/").
        Handler(http.StripPrefix("/static/style/", http.FileServer(http.Dir("../static/style/"))))

    log.Println("Routers created")
    return router
}
