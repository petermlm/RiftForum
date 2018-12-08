package main

import (
    "math/rand"
)

func make_new_invite() *Invite {
    // TODO: This is horrible
    var letter_runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    db := GetDBCon()
    key := make([]rune, 64)

    for i:=0; i<64; i++ {
        rn := rand.Intn(len(letter_runes))
        key[i] = letter_runes[rn]
    }

    key_str := string(key)

    invite := &Invite{
        Key: key_str,
        Status: Unused,
    }

    err := db.Insert(invite)

    if err != nil {
        panic(err)
    }

    return invite
}
