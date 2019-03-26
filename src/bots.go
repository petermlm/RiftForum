package main

import "fmt"

func BasicBot() {
    fmt.Println("Running")
}

func HelloBot() {
    user := get_rift_bot_user()
    NewTopic(user, "I'm a Bot", "And this is a message written by a Bot")
}

func get_rift_bot_user() *User {
    user := new(User)

    err := db.Model(user).Where("Username = 'RiftBot'").Select()

    if err != nil {
        panic(err)
    }

    return user
}
