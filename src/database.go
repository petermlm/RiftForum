package main

import (
    "log"

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

    rift_bot := NewUser("RiftBot", Basic, "pl")
    err = db.Insert(rift_bot)

    if err != nil {
        return err
    }

    return nil
}
