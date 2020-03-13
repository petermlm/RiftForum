package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("Creating default users...")

		// Admin and first user
		NewUser(AdminUsername, Administrator, DefaultPassword)
		NewUser(FirstUsername, Moderator, DefaultPassword)

		// Bots
		NewUser("RiftBot", Bot, DefaultPassword)
		NewUser("GreeterBot", Bot, DefaultPassword)
		NewUser("RedditBot", Bot, DefaultPassword)
		NewUser("YoutubeBot", Bot, DefaultPassword)

		return nil
	}, func(db migrations.DB) error {
		fmt.Println("Dropping default my_table...")
		_, err := db.Exec("TRUNCATE Users")
		return err
	})
}
