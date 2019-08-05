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
    InviteSet(invite_key, Used)
}

func ChangePassword(user *User, new_password, new_password2 string) {
    if new_password != new_password2 {
        return
    }

    hash, err := GenerateHash(new_password)

    if err != nil {
        panic(err)
    }

    db := GetDBCon()
    user.PasswordHash = hash
    db.Update(user)
}
