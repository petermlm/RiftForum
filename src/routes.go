package main

import (
    "context"
    "fmt"
    "log"
    "strconv"
    "net/http"

    "github.com/gorilla/mux"
)

func index(res http.ResponseWriter, req *http.Request) {
    topics := GetTopics()
    data := SerializeTopics(topics)
    Render(&res, req, "index.html", data)
}

func login_get(res http.ResponseWriter, req *http.Request) {
    ctx := req.Context()
    user_info := ctx.Value("UserInfo")

    if user_info != nil {
        Redirect(&res, req, "/")
        return
    }

    data := SerializeEmpty()
    Render(&res, req, "login.html", data)
}

func login_post(res http.ResponseWriter, req *http.Request) {
    form_username := req.PostFormValue("username")
    form_password := req.PostFormValue("password")

    token, err := CreateToken(form_username, form_password)

    if err == nil {
        cookie := http.Cookie{
            Name: "jwt",
            Value: token,
        }
        http.SetCookie(res, &cookie)
    }

    Redirect(&res, req, "/")
}

func logout(res http.ResponseWriter, req *http.Request) {
    cookie := http.Cookie{
        Name: "jwt",
        Value: "",
    }
    http.SetCookie(res, &cookie)

    Redirect(&res, req, "/")
}

func register_get(res http.ResponseWriter, req *http.Request) {
    key := req.URL.Query().Get("key")
    data := SerializeRegister(key)
    Render(&res, req, "register.html", data)
}

func register_post(res http.ResponseWriter, req *http.Request) {
    form_invite_key := req.PostFormValue("invite_key")
    form_username := req.PostFormValue("username")
    form_password := req.PostFormValue("password")
    form_password2 := req.PostFormValue("password2")

    Register(form_invite_key, form_username, form_password, form_password2)

    Redirect(&res, req, "/")
}

func admin_get(res http.ResponseWriter, req *http.Request) {
    data := SerializeEmpty()
    Render(&res, req, "admin.html", data)
}

func admin_invites_get(res http.ResponseWriter, req *http.Request) {
    invites := GetInvites()
    data := SerializeInvites(invites)
    Render(&res, req, "invites.html", data)
}

func admin_invites_post(res http.ResponseWriter, req *http.Request) {
    new_invite := make_new_invite()
    data := SerializeInviteNew(new_invite)
    Render(&res, req, "invite_new.html", data)
}

func admin_invites_cancel_get(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    key := vars["key"]

    InviteSet(key, Canceled)

    Redirect(&res, req, "/admin/invites")
}

func admin_cancel_all_post(res http.ResponseWriter, req *http.Request) {
    InviteCancelAll()
    Redirect(&res, req, "/admin/invites")
}

func admin_users_change_type_get(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    username := vars["username"]
    new_type_str := req.URL.Query().Get("new_type")

    var new_type UserTypes
    set := false

    if new_type_str == "basic" {
        new_type = Basic
        set = true
    } else if new_type_str == "moderator" {
        new_type = Moderator
        set = true
    }

    if set {
        UserTypeSet(username, new_type)
    }

    Redirect(&res, req, "/users")
}

func topics_post(res http.ResponseWriter, req *http.Request) {
    var err error
    db := GetDBCon()
    form_title := req.PostFormValue("title")
    form_message := req.PostFormValue("message")

    // Author
    user := new(User)
    ctx := req.Context()
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

    Redirect(&res, req, "/")
}

func topic_get(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    topic_id_parsed, err := strconv.ParseUint(vars["id"], 10, 32)

    if err != nil {
        NotFound(&res, req)
        return
    }

    topic_id := uint(topic_id_parsed)
    topic := GetTopic(topic_id)

    if topic == nil {
        NotFound(&res, req)
        return
    }

    data := SerializeTopic(topic)
    Render(&res, req, "topic.html", data)
}

func topic_post(res http.ResponseWriter, req *http.Request) {
    var err error
    db := GetDBCon()
    form_message := req.PostFormValue("message")

    // Author
    user := new(User)
    ctx := req.Context()
    user_info := ctx.Value("UserInfo").(*UserInfo)
    err = db.Model(user).Where("Username = ?", user_info.Username).Select()

    if err != nil {
        panic(err)
    }

    // Topic
    vars := mux.Vars(req)
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
    Redirect(&res, req, redirect_path)
}

func users_get(res http.ResponseWriter, req *http.Request) {
    users := GetUsers()
    data := SerializeUsers(users)
    Render(&res, req, "users.html", data)
}

func user_get(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    username := vars["username"]

    user, err := GetUser(username)

    if err != nil {
        return
    }

    data := SerializeUser(user)
    Render(&res, req, "user.html", data)
}

func user_about_post(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    username := vars["username"]
    form_new_about := req.PostFormValue("about")

    UserSetAbout(username, form_new_about)
    redirect_path := fmt.Sprintf("/users/%s", username)
    Redirect(&res, req, redirect_path)
}

func user_signature_post(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    username := vars["username"]
    form_new_signature := req.PostFormValue("signature")

    UserSetSignature(username, form_new_signature)
    redirect_path := fmt.Sprintf("/users/%s", username)
    Redirect(&res, req, redirect_path)
}

func save_user_info(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        cookie, err := req.Cookie("jwt")
        ctx := req.Context()
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

        next.ServeHTTP(w, req.WithContext(ctx))
    })
}

func auth_middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        ctx := req.Context()
        user_info_value := ctx.Value("UserInfo")
        is_auth := false

        if user_info_value != nil {
            user_info := user_info_value.(*UserInfo)

            if user_info != nil {
                is_auth = true
            }
        }

        if is_auth {
            next.ServeHTTP(w, req.WithContext(ctx))
        } else {
            http.Error(w, "Forbidden", http.StatusForbidden)
        }
    })
}

func admin_middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        ctx := req.Context()
        user_info := ctx.Value("UserInfo").(*UserInfo)

        if user_info != nil {
            next.ServeHTTP(w, req.WithContext(ctx))
        } else {
            http.Error(w, "Forbidden", http.StatusForbidden)
        }
    })
}

func CreateRouter() *mux.Router {
    router := mux.NewRouter()
    router.HandleFunc("/", index).Methods("GET")
    router.HandleFunc("/login", login_get).Methods("GET")
    router.HandleFunc("/login", login_post).Methods("POST")
    router.HandleFunc("/logout", logout).Methods("POST")
    router.HandleFunc("/register", register_get).Methods("GET")
    router.HandleFunc("/register", register_post).Methods("POST")
    router.Use(save_user_info)

    auth_routes := router.PathPrefix("/").Subrouter()
    auth_routes.HandleFunc("/topics", topics_post).Methods("POST")
    auth_routes.HandleFunc("/topics/{id:[0-9]+}", topic_get).Methods("GET")
    auth_routes.HandleFunc("/topics/{id:[0-9]+}", topic_post).Methods("POST")
    auth_routes.HandleFunc("/users", users_get).Methods("GET")
    auth_routes.HandleFunc("/users/{username:[a-zA-Z0-9]+}", user_get).Methods("GET")
    auth_routes.HandleFunc("/users/{username:[a-zA-Z0-9]+}/about", user_about_post).Methods("POST")
    auth_routes.HandleFunc("/users/{username:[a-zA-Z0-9]+}/signature", user_signature_post).Methods("POST")
    auth_routes.Use(auth_middleware)

    admin_routes := auth_routes.PathPrefix("/admin").Subrouter()
    admin_routes.HandleFunc("/", admin_get).Methods("GET")
    admin_routes.HandleFunc("/invites", admin_invites_get).Methods("GET")
    admin_routes.HandleFunc("/invites_new", admin_invites_post).Methods("GET")
    admin_routes.HandleFunc("/invites_cancel/{key:[a-zA-Z0-9]+}", admin_invites_cancel_get).Methods("GET")
    admin_routes.HandleFunc("/invites_cancel_all", admin_cancel_all_post).Methods("GET")
    admin_routes.HandleFunc("/users/{username:[a-zA-Z0-9]+}/change_type", admin_users_change_type_get).Methods("GET")
    admin_routes.Use(admin_middleware)

    router.
        PathPrefix("/static/style/").
        Handler(http.StripPrefix("/static/style/", http.FileServer(http.Dir("../static/style/"))))

    log.Println("Routers created")
    return router
}
