package main

func Register(invite_key string,
              username string,
              password string,
              password2 string) {
    if password != password2 {
        return
    }

    if !InviteExists(invite_key) {
        return
    }

    new_user := NewUser(username, Basic, password)
    SaveUser(new_user)
}
