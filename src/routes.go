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
    Render(&writer, r, "index.html", data)
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
    Render(&writer, r, "topic.html", data)
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

func save_user_info(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("jwt")
        ctx := r.Context()
        authenticated := false

        if err == nil {
            claims := VerifyToken(cookie.Value)

            if claims != nil {
                user_info := &UserInfo {
                    Username: claims.Username,
                }

                ctx = context.WithValue(ctx, "UserInfo", user_info)
                authenticated = true
            }
        }

        if !authenticated {
            ctx = context.WithValue(ctx, "UserInfo", nil)
        }

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func auth_middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        user_info := ctx.Value("UserInfo").(*UserInfo)

        if user_info != nil {
            next.ServeHTTP(w, r.WithContext(ctx))
        } else {
            http.Error(w, "Forbidden", http.StatusForbidden)
        }
    })
}

func CreateRouter() *mux.Router {
    router := mux.NewRouter()
    router.HandleFunc("/", index).Methods("GET")
    router.HandleFunc("/login", login).Methods("POST")
    router.HandleFunc("/logout", logout).Methods("POST")
    router.Use(save_user_info)

    auth_routes := router.PathPrefix("/").Subrouter()
    auth_routes.HandleFunc("/topics", topics_post).Methods("POST")
    auth_routes.HandleFunc("/topics/{id:[0-9]+}", topic_get).Methods("GET")
    auth_routes.HandleFunc("/topics/{id:[0-9]+}", topic_post).Methods("POST")
    auth_routes.Use(auth_middleware)

    router.
        PathPrefix("/static/style/").
        Handler(http.StripPrefix("/static/style/", http.FileServer(http.Dir("../static/style/"))))

    log.Println("Routers created")
    return router
}
