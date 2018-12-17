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

func register_get(writer http.ResponseWriter, r *http.Request) {
    key := r.URL.Query().Get("key")
    data := SerializeRegister(key)
    Render(&writer, r, "register.html", data)
}

func register_post(writer http.ResponseWriter, r *http.Request) {
    form_invite_key := r.PostFormValue("invite_key")
    form_username := r.PostFormValue("username")
    form_password := r.PostFormValue("password")
    form_password2 := r.PostFormValue("password2")

    Register(form_invite_key, form_username, form_password, form_password2)

    http.Redirect(writer, r, "/", http.StatusSeeOther)
}

func admin_get(writer http.ResponseWriter, r *http.Request) {
    data := SerializeEmpty()
    Render(&writer, r, "admin.html", data)
}

func admin_invites_get(writer http.ResponseWriter, r *http.Request) {
    invites := GetInvites()
    data := SerializeInvites(invites)
    Render(&writer, r, "invites.html", data)
}

func admin_invites_post(writer http.ResponseWriter, r *http.Request) {
    new_invite := make_new_invite()
    data := SerializeInviteNew(new_invite)
    Render(&writer, r, "invite_new.html", data)
}

func admin_invites_cancel_get(writer http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    key := vars["key"]

    InviteSet(key, Canceled)

    http.Redirect(writer, r, "/admin/invites", http.StatusSeeOther)
}

func admin_cancel_all_post(writer http.ResponseWriter, r *http.Request) {
    InviteCancelAll()
    http.Redirect(writer, r, "/admin/invites", http.StatusSeeOther)
}

func topics_post(writer http.ResponseWriter, r *http.Request) {
    var err error
    db := GetDBCon()
    form_title := r.PostFormValue("title")
    form_message := r.PostFormValue("message")

    // Author
    user := new(User)
    ctx := r.Context()
    user_info := ctx.Value("UserInfo").(*UserInfo)
    err = db.Model(user).Where("Username = ?", user_info.Username).Select()

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
    ctx := r.Context()
    user_info := ctx.Value("UserInfo").(*UserInfo)
    err = db.Model(user).Where("Username = ?", user_info.Username).Select()

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

func users_get(writer http.ResponseWriter, r *http.Request) {
    users := GetUsers()
    data := SerializeUsers(users)
    Render(&writer, r, "users.html", data)
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
                    Usertype: claims.Usertype,
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
        user_info_value := ctx.Value("UserInfo")
        is_auth := false

        if user_info_value != nil {
            user_info := user_info_value.(*UserInfo)

            if user_info != nil {
                is_auth = true
            }
        }

        if is_auth {
            next.ServeHTTP(w, r.WithContext(ctx))
        } else {
            http.Error(w, "Forbidden", http.StatusForbidden)
        }
    })
}

func admin_middleware(next http.Handler) http.Handler {
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
    router.HandleFunc("/register", register_get).Methods("GET")
    router.HandleFunc("/register", register_post).Methods("POST")
    router.Use(save_user_info)

    auth_routes := router.PathPrefix("/").Subrouter()
    auth_routes.HandleFunc("/topics", topics_post).Methods("POST")
    auth_routes.HandleFunc("/topics/{id:[0-9]+}", topic_get).Methods("GET")
    auth_routes.HandleFunc("/topics/{id:[0-9]+}", topic_post).Methods("POST")
    auth_routes.HandleFunc("/users", users_get).Methods("GET")
    auth_routes.Use(auth_middleware)

    admin_routes := auth_routes.PathPrefix("/admin").Subrouter()
    admin_routes.HandleFunc("/", admin_get).Methods("GET")
    admin_routes.HandleFunc("/invites", admin_invites_get).Methods("GET")
    admin_routes.HandleFunc("/invites_new", admin_invites_post).Methods("GET")
    admin_routes.HandleFunc("/invites_cancel/{key:[a-zA-Z0-9]+}", admin_invites_cancel_get).Methods("GET")
    admin_routes.HandleFunc("/invites_cancel_all", admin_cancel_all_post).Methods("GET")
    admin_routes.Use(admin_middleware)

    router.
        PathPrefix("/static/style/").
        Handler(http.StripPrefix("/static/style/", http.FileServer(http.Dir("../static/style/"))))

    log.Println("Routers created")
    return router
}
