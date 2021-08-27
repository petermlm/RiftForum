package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("Creating default users...")

		// Admin and first user
		NewUser(Config.AdminUsername, Administrator, Config.DefaultPassword)
		NewUser(Config.FirstUsername, Moderator, Config.DefaultPassword)

		// Bots
		NewUser("RiftBot", Bot, Config.DefaultPassword)
		NewUser("GreeterBot", Bot, Config.DefaultPassword)
		NewUser("RedditBot", Bot, Config.DefaultPassword)
		NewUser("YoutubeBot", Bot, Config.DefaultPassword)

		return nil
	}, func(db migrations.DB) error {
		fmt.Println("Dropping default my_table...")
		_, err := db.Exec("TRUNCATE Users")
		return err
	})
}
