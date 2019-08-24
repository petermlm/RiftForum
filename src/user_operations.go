package main

import "errors"

func Register(invite_key string,
              username string,
              password string,
              password2 string) error {
    if password != password2 {
        return errors.New("Passwords don't match")
    }

    if !InviteExists(invite_key) {
        return errors.New("Invite doesn't exist")
    }

    if len(username) > MaxUsernameSize {
        return errors.New("Username is to big")
    }

    NewUser(username, Basic, password)
    InviteSet(invite_key, Used)

    return nil
}

func ChangePassword(user *User, new_password, new_password2 string) error {
    if new_password != new_password2 {
        return errors.New("Passwords don't match")
    }

    hash, err := GenerateHash(new_password)

    if err != nil {
        RiftForumPanic("Can't generate hash for password", err)
    }

    db := GetDBCon()
    user.PasswordHash = hash
    db.Update(user)

    return nil
}
