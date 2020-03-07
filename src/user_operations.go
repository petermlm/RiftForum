package main

import "errors"

func Register(invite_key string,
              username string,
              password string,
              password2 string) *RegisterErrors {
    register_errors := &RegisterErrors{false, false, false, false}

    if password != password2 {
        register_errors.passwords_dont_match = true
        return register_errors
    }

    if !InviteExists(invite_key) {
        register_errors.invite_key_bad = true
        return register_errors
    }

    if len(username) > MaxUsernameSize {
        register_errors.username_is_invalid = true
        return register_errors
    }

    if _, err := GetUser(username); err == nil {
        register_errors.username_alreay_taken = true
        return register_errors
    }

	_, err := NewUser(username, Basic, password)
    if err != nil {
        register_errors.username_is_invalid = true
        return register_errors
    }

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
