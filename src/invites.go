package main

import (
	"math/rand"
)

var letter_runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func make_new_invite() *Invite {
	db := GetDBCon()
	key := make([]rune, Config.InviteSize)

	for i := 0; i < Config.InviteSize; i++ {
		key[i] = letter_runes[rand.Intn(len(letter_runes))]
	}

	key_str := string(key)

	invite := &Invite{
		Key:    key_str,
		Status: Unused,
	}

	err := db.Insert(invite)

	if err != nil {
		panic(err)
	}

	return invite
}
