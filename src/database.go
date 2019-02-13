package main

import (
    "log"
    "fmt"

    "github.com/go-pg/pg"
    "github.com/go-pg/pg/orm"
)

var db *pg.DB

func InitDB() {
    var err error

    db = pg.Connect(&pg.Options{
        Addr: "postgres:5432",
        Database: "riftforum_db",
        User: "riftforum_user",
        Password: "riftforum_pass",

    })

    err = createSchema()

    if err != nil {
        panic(err)
    }

    err = createDefaultData()

    if err != nil {
        panic(err)
    }

    log.Println("Database connection astablished")
}

func CloseDB() {
    db.Close()
    log.Println("Database connection closed")
}

func GetDBCon() *pg.DB {
    return db
}

func createSchema() error {
    for _, model := range []interface{}{
        (*User)(nil),
        (*Topic)(nil),
        (*Message)(nil),
        (*Invite)(nil),
    } {
        err := db.CreateTable(model, &orm.CreateTableOptions{
            Temp: true,
        })

        if err != nil {
            return err
        }
    }

    return nil
}

func createDefaultData() error {
    var err error

    // Users
    user_admin := NewUser("admin", Administrator, "pl")
    err = db.Insert(user_admin)

    if err != nil {
        return err
    }

    user_basic := NewUser("petermlm", Moderator, "pl")
    err = db.Insert(user_basic)

    if err != nil {
        return err
    }

    // Topic
    for i:=0; i<100; i++ {
        title := fmt.Sprintf("Test Topic %d", i)

        topic := &Topic{
            Title: title,
            Author: user_admin,
            AuthorId: user_admin.Id,
        }

        err = db.Insert(topic)

        if err != nil {
            return err
        }

        // Messages
        if i == 99 {
            for j:=0; j<100; j++ {
                message := &Message{
                    Message: fmt.Sprintf("Test message %d", j),
                    Author: user_admin,
                    AuthorId: user_admin.Id,
                    Topic: topic,
                    TopicId: topic.Id,
                }

                err = db.Insert(message)
            }
        } else {
            message := &Message{
                Message: "Test message",
                Author: user_admin,
                AuthorId: user_admin.Id,
                Topic: topic,
                TopicId: topic.Id,
            }

            err = db.Insert(message)
        }

        if err != nil {
            return err
        }
    }

    return nil
}
