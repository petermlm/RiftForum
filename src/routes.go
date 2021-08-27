package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func index(res http.ResponseWriter, req *http.Request) {
	page := PageFromRequest(req)
	topics := GetTopics(page)
	data := SerializeTopics(topics, page)
	Render(&res, req, "index.html", data)
}

func login_get(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	user_info := ctx.Value("UserInfo")

	if user_info != nil {
		Redirect(&res, req, "/")
		return
	}

	next_page := req.URL.Query().Get("next_page")
	data := SerializeLogin(next_page)
	Render(&res, req, "login.html", data)
}

func login_post(res http.ResponseWriter, req *http.Request) {
	form_username := req.PostFormValue("username")
	form_password := req.PostFormValue("password")
	next_page := req.PostFormValue("next_page")

	token, err := CreateToken(form_username, form_password)

	if err == nil {
		cookie := http.Cookie{
			Name:  "jwt",
			Value: token,
		}
		http.SetCookie(res, &cookie)
	}

	if next_page == "" {
		next_page = "/"
	}

	Redirect(&res, req, next_page)
}

func logout_post(res http.ResponseWriter, req *http.Request) {
	cookie := http.Cookie{
		Name:  "jwt",
		Value: "",
	}
	http.SetCookie(res, &cookie)
	Redirect(&res, req, "/")
}

func register_get(res http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")
	errors := RegisterErrors{false, false, false, false}
	data := SerializeRegister(key, errors)
	Render(&res, req, "register.html", data)
}

func register_post(res http.ResponseWriter, req *http.Request) {
	form_invite_key := req.PostFormValue("invite_key")
	form_username := req.PostFormValue("username")
	form_password := req.PostFormValue("password")
	form_password2 := req.PostFormValue("password2")

	errors := Register(form_invite_key, form_username, form_password, form_password2)

	if errors == nil {
		Redirect(&res, req, "/")
	} else {
		data := SerializeRegister(form_invite_key, *errors)
		Render(&res, req, "register.html", data)
	}
}

func admin_get(res http.ResponseWriter, req *http.Request) {
	data := SerializeEmpty()
	Render(&res, req, "admin.html", data)
}

func admin_invites_get(res http.ResponseWriter, req *http.Request) {
	page := PageFromRequest(req)
	invites := GetInvites(page)
	data := SerializeInvites(invites, page)
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

func admin_users_change_password_get(res http.ResponseWriter, req *http.Request) {
	errors := ChangePasswordErrors{false, false}
	render_change_password(res, req, true, errors)
}

func admin_users_change_password_post(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	form_new_password := req.PostFormValue("new_password")
	form_new_password2 := req.PostFormValue("new_password2")

	user, err := GetUser(username)

	if err != nil {
		NotFound(&res, req)
		return
	}

	err = ChangePassword(user, form_new_password, form_new_password2)

	if err != nil {
		errors := ChangePasswordErrors{false, true}
		render_change_password(res, req, false, errors)
		return
	}

	Redirect(&res, req, fmt.Sprintf("/users/%s", username))
}

func admin_users_ban_get(res http.ResponseWriter, req *http.Request) {
	user, err := get_user_from_url(req)

	if err != nil {
		NotFound(&res, req)
		return
	}

	BanUser(user)
	Redirect(&res, req, fmt.Sprintf("/users/%s", user.Username))
}

func admin_users_unban_get(res http.ResponseWriter, req *http.Request) {
	user, err := get_user_from_url(req)

	if err != nil {
		NotFound(&res, req)
		return
	}

	UnbanUser(user)
	Redirect(&res, req, fmt.Sprintf("/users/%s", user.Username))
}

func topics_post(res http.ResponseWriter, req *http.Request) {
	var err error
	db := GetDBCon()
	form_title := req.PostFormValue("title")
	form_message := req.PostFormValue("message")

	form_title = strings.TrimSpace(form_title)
	form_title = strings.ReplaceAll(form_title, "\r\n", " ")
	form_message = strings.TrimSpace(form_message)

	if len(form_title) == 0 || len(form_message) == 0 {
		OperationNotAllowed(&res, req)
		return
	}

	// Author
	user := new(User)
	ctx := req.Context()
	user_info := ctx.Value("UserInfo").(*UserInfo)
	err = db.Model(user).Where("Username = ?", user_info.Username).Select()

	if err != nil {
		panic(err)
	}

	topic := NewTopic(user, form_title, form_message)
	Redirect(&res, req, fmt.Sprintf("/topics/%d", topic.Id))
}

func topic_get(res http.ResponseWriter, req *http.Request) {
	page := PageFromRequest(req)
	vars := mux.Vars(req)
	topic_id_parsed, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		NotFound(&res, req)
		return
	}

	topic_id := uint(topic_id_parsed)
	topic := GetTopic(topic_id, page)

	if topic == nil {
		NotFound(&res, req)
		return
	}

	data := SerializeTopic(topic, page)
	Render(&res, req, "topic.html", data)
}

func topic_post(res http.ResponseWriter, req *http.Request) {
	var err error
	db := GetDBCon()
	form_message := req.PostFormValue("message")

	if len(form_message) == 0 {
		OperationNotAllowed(&res, req)
		return
	}

	// Author
	user := new(User)
	ctx := req.Context()
	user_info := ctx.Value("UserInfo").(*UserInfo)
	err = db.Model(user).Where("Username = ?", user_info.Username).Select()

	if err != nil {
		RiftForumPanic("Could not find user", err)
	}

	// Topic
	vars := mux.Vars(req)
	topic_id_parsed, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		NotFound(&res, req)
		return
	}
	topic_id := uint(topic_id_parsed)
	topic := GetTopic(topic_id, PageDefault())

	// Message
	NewMessage(user, topic, form_message)

	// Get page to redirect to
	var page Page
	var page_querystr string

	msg_pages := CountMessagePages(topic_id, PageDefault())
	if msg_pages == 0 {
		page = PageDefault()
	} else {
		page = NewPage(msg_pages)
		page_querystr = "?" + page.querystr()
	}

	redirect_path := fmt.Sprintf("/topics/%d%s", topic_id, page_querystr)
	Redirect(&res, req, redirect_path)
}

func message_get(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	user_info := ctx.Value("UserInfo").(*UserInfo)
	vars := mux.Vars(req)
	message_id_parsed, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		NotFound(&res, req)
		return
	}

	message_id := uint(message_id_parsed)
	message := GetMessage(message_id)

	if message == nil {
		NotFound(&res, req)
		return
	}

	if !check_permission(message.Author.Username, user_info) {
		OperationNotAllowed(&res, req)
		return
	}

	data := SerializeMessageEdit(message)
	Render(&res, req, "message_edit.html", data)
}

func message_post(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	user_info := ctx.Value("UserInfo").(*UserInfo)
	vars := mux.Vars(req)
	message_id_parsed, err := strconv.ParseUint(vars["id"], 10, 32)
	form_message := req.PostFormValue("message")

	if err != nil {
		NotFound(&res, req)
		return
	}

	message_id := uint(message_id_parsed)
	message := GetMessage(message_id)
	message.Message = form_message
	UpdateMessage(message)

	if message == nil {
		NotFound(&res, req)
		return
	}

	if !check_permission(message.Author.Username, user_info) {
		OperationNotAllowed(&res, req)
		return
	}

	redirect_path := fmt.Sprintf("/topics/%d", message.Topic.Id)
	Redirect(&res, req, redirect_path)
}

func users_get(res http.ResponseWriter, req *http.Request) {
	page := PageFromRequest(req)
	users := GetUsers(page)
	data := SerializeUsers(users, page)
	Render(&res, req, "users.html", data)
}

func user_get(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	user_info := ctx.Value("UserInfo").(*UserInfo)

	vars := mux.Vars(req)
	username := vars["username"]

	user, err := GetUser(username)

	if err != nil {
		NotFound(&res, req)
		return
	}

	data := SerializeUser(user_info.Username, user)
	Render(&res, req, "user.html", data)
}

func user_about_post(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	user_info := ctx.Value("UserInfo").(*UserInfo)

	vars := mux.Vars(req)
	username := vars["username"]
	form_new_about := req.PostFormValue("about")

	if !check_permission(username, user_info) {
		OperationNotAllowed(&res, req)
		return
	}

	UserSetAbout(username, form_new_about)
	redirect_path := fmt.Sprintf("/users/%s", username)
	Redirect(&res, req, redirect_path)
}

func user_signature_post(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	user_info := ctx.Value("UserInfo").(*UserInfo)

	vars := mux.Vars(req)
	username := vars["username"]
	form_new_signature := req.PostFormValue("signature")

	if !check_permission(username, user_info) {
		OperationNotAllowed(&res, req)
		return
	}

	UserSetSignature(username, form_new_signature)
	redirect_path := fmt.Sprintf("/users/%s", username)
	Redirect(&res, req, redirect_path)
}

func user_change_password_get(res http.ResponseWriter, req *http.Request) {
	errors := ChangePasswordErrors{false, false}
	render_change_password(res, req, false, errors)
}

func user_change_password_post(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	form_old_password := req.PostFormValue("old_password")
	form_new_password := req.PostFormValue("new_password")
	form_new_password2 := req.PostFormValue("new_password2")

	user, err := GetUser(username)

	if err != nil {
		NotFound(&res, req)
		return
	}

	if !VerifyUserPass(user, form_old_password) {
		errors := ChangePasswordErrors{true, false}
		render_change_password(res, req, false, errors)
		return
	}

	err = ChangePassword(user, form_new_password, form_new_password2)

	if err != nil {
		errors := ChangePasswordErrors{false, true}
		render_change_password(res, req, false, errors)
		return
	}

	Redirect(&res, req, fmt.Sprintf("/users/%s", username))
}

func bots_get(res http.ResponseWriter, req *http.Request) {
	hearthbeat_status := GetHearthBeatStatus()
	data := SerializeBots(hearthbeat_status)
	Render(&res, req, "bots.html", data)
}

func bbcode_get(res http.ResponseWriter, req *http.Request) {
	data := SerializeEmpty()
	Render(&res, req, "bbcode.html", data)
}

func save_user_info(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("jwt")
		ctx := req.Context()
		authenticated := false

		if err == nil {
			claims := VerifyToken(cookie.Value)

			if claims != nil {
				user, _ := GetUser(claims.Username)

				if !user.Banned {
					user_info := &UserInfo{
						Id:       claims.Id,
						Username: claims.Username,
						Usertype: claims.Usertype,
					}

					ctx = context.WithValue(ctx, "UserInfo", user_info)
					authenticated = true
				}
			}
		}

		if !authenticated {
			ctx = context.WithValue(ctx, "UserInfo", nil)
		}

		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

func auth_middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
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
			next.ServeHTTP(res, req.WithContext(ctx))
		} else {
			Login(&res, req, req.URL.Path)
		}
	})
}

func admin_middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user_info := ctx.Value("UserInfo").(*UserInfo)

		if user_info != nil && user_info.Usertype == Administrator {
			next.ServeHTTP(res, req.WithContext(ctx))
		} else {
			AdminOnly(&res, req)
		}
	})
}

func get_key_from_url(r *http.Request, key string) (string, error) {
	vars := mux.Vars(r)
	if value, ok := vars[key]; ok {
		return value, nil
	}

	err_str := fmt.Sprintf("Not in request: %s", key)
	return "", errors.New(err_str)
}

func get_user_from_url(r *http.Request) (*User, error) {
	username, err := get_key_from_url(r, "username")
	if err != nil {
		return nil, err
	}
	return GetUser(username)
}

func check_permission(username string, user_info *UserInfo) bool {
	return username == user_info.Username || user_info.IsMod()
}

func render_change_password(
	res http.ResponseWriter,
	req *http.Request,
	is_for_admin bool,
	errors ChangePasswordErrors,
) {
	vars := mux.Vars(req)
	username := vars["username"]

	user, err := GetUser(username)

	if err != nil {
		NotFound(&res, req)
		return
	}

	data := SerializeChangePassword(user, is_for_admin, errors)
	Render(&res, req, "change_password.html", data)
}

func CreateRouter() *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/login", login_get).Methods("GET")
	router.HandleFunc("/login", login_post).Methods("POST")
	router.HandleFunc("/logout", logout_post).Methods("POST")
	router.HandleFunc("/register", register_get).Methods("GET")
	router.HandleFunc("/register", register_post).Methods("POST")
	router.Use(save_user_info)

	auth_routes := router.PathPrefix("/").Subrouter()
	auth_routes.HandleFunc("/topics", topics_post).Methods("POST")
	auth_routes.HandleFunc("/topics/{id:[0-9]+}", topic_get).Methods("GET")
	auth_routes.HandleFunc("/topics/{id:[0-9]+}", topic_post).Methods("POST")
	auth_routes.HandleFunc("/messages/{id:[0-9]+}", message_get).Methods("GET")
	auth_routes.HandleFunc("/messages/{id:[0-9]+}", message_post).Methods("POST")
	auth_routes.HandleFunc("/users", users_get).Methods("GET")
	auth_routes.HandleFunc("/users/{username:[a-zA-Z0-9]+}", user_get).Methods("GET")
	auth_routes.HandleFunc("/users/{username:[a-zA-Z0-9]+}/about", user_about_post).Methods("POST")
	auth_routes.HandleFunc("/users/{username:[a-zA-Z0-9]+}/signature", user_signature_post).Methods("POST")
	auth_routes.HandleFunc("/users/{username:[a-zA-Z0-9]+}/change_password", user_change_password_get).Methods("GET")
	auth_routes.HandleFunc("/users/{username:[a-zA-Z0-9]+}/change_password", user_change_password_post).Methods("POST")
	auth_routes.HandleFunc("/bots", bots_get).Methods("GET")
	auth_routes.HandleFunc("/bbcode", bbcode_get).Methods("GET")
	auth_routes.Use(auth_middleware)

	admin_routes := auth_routes.PathPrefix("/admin").Subrouter()
	admin_routes.HandleFunc("/", admin_get).Methods("GET")
	admin_routes.HandleFunc("/invites", admin_invites_get).Methods("GET")
	admin_routes.HandleFunc("/invites_new", admin_invites_post).Methods("GET")
	admin_routes.HandleFunc("/invites_cancel/{key:[a-zA-Z0-9]+}", admin_invites_cancel_get).Methods("GET")
	admin_routes.HandleFunc("/invites_cancel_all", admin_cancel_all_post).Methods("GET")
	admin_routes.HandleFunc("/users/{username:[a-zA-Z0-9]+}/change_type", admin_users_change_type_get).Methods("GET")
	admin_routes.HandleFunc("/users/{username:[a-zA-Z0-9]+}/change_password", admin_users_change_password_get).Methods("GET")
	admin_routes.HandleFunc("/users/{username:[a-zA-Z0-9]+}/change_password", admin_users_change_password_post).Methods("POST")
	admin_routes.HandleFunc("/users/{username:[a-zA-Z0-9]+}/ban", admin_users_ban_get).Methods("GET")
	admin_routes.HandleFunc("/users/{username:[a-zA-Z0-9]+}/unban", admin_users_unban_get).Methods("GET")
	admin_routes.Use(admin_middleware)

	router.
		PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	log.Println("Routers created")
	return router
}
