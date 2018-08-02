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
    err := db.Insert(&User{
        Username: "admin",
        // PasswordHash: "",
        Signature: "I'm the Administrator",
        About: "I'm the Administrator",
        UserType: Administrator,
    })

    if err != nil {
        return err
    }

    return nil
}
