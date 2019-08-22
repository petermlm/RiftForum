package main

import (
    "fmt"

    "github.com/go-pg/migrations"
)

func init() {
    migrations.MustRegisterTx(func(db migrations.DB) error {
        fmt.Println("Creating default users...")
        NewUser("admin", Administrator, "pl")
        NewUser("petermlm", Moderator, "pl")
        NewUser("RiftBot", Bot, "pl")
        NewUser("GreeterBot", Bot, "pl")
        NewUser("RedditBot", Bot, "pl")
        NewUser("YoutubeBot", Bot, "pl")
        return nil
    }, func(db migrations.DB) error {
        fmt.Println("Dropping default my_table...")
        _, err := db.Exec("TRUNCATE Users")
        return err
    })
}
