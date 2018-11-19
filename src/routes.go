package main

import (
    "context"
    "fmt"
    "log"
    "strconv"
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

    ctx := r.Context()
    fmt.Println(ctx.Value("Username"))

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

func login(writer http.ResponseWriter, r *http.Request) {
    form_username := r.PostFormValue("username")
    form_password := r.PostFormValue("password")

    token, err := CreateToken(form_username, form_password)

    if err == nil {
        cookie := http.Cookie{
            Name: "jwt",
            Value: token,
        }
        http.SetCookie(writer, &cookie)
    }

    http.Redirect(writer, r, "/", http.StatusSeeOther)
}

func logout(writer http.ResponseWriter, r *http.Request) {
    cookie := http.Cookie{
        Name: "jwt",
        Value: "",
    }
    http.SetCookie(writer, &cookie)

    http.Redirect(writer, r, "/", http.StatusSeeOther)
}

func auth_middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var authenticated bool
        cookie, ok := r.Cookie("jwt")

        if ok == nil {
            valid := VerifyToken(cookie.Value)

            if valid {
                ctx := context.WithValue(r.Context(), "Username", "ze")
                next.ServeHTTP(w, r.WithContext(ctx))
                authenticated = true
            } else {
                authenticated = false
            }
        } else {
            authenticated = false
        }

        if !authenticated {
            http.Error(w, "Forbidden", http.StatusForbidden)
        }
    })
}

func CreateRouter() *mux.Router {
    router := mux.NewRouter()
    router.HandleFunc("/", index).Methods("GET")
    router.HandleFunc("/login", login).Methods("POST")
    router.HandleFunc("/logout", logout).Methods("POST")

    topics_router := router.PathPrefix("/topics").Subrouter()
    topics_router.HandleFunc("/", topics_post).Methods("POST")
    topics_router.HandleFunc("/{id:[0-9]+}", topic_get).Methods("GET")
    topics_router.HandleFunc("/{id:[0-9]+}", topic_post).Methods("POST")
    topics_router.Use(auth_middleware)

    router.
        PathPrefix("/static/style/").
        Handler(http.StripPrefix("/static/style/", http.FileServer(http.Dir("../static/style/"))))

    log.Println("Routers created")
    return router
}
