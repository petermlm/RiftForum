package main

import (
    "time"
    "log"

    "github.com/go-pg/pg"
    "github.com/go-pg/pg/orm"
)

var db *pg.DB

func InitDB() {
    // Connect to database, waits one second between attempts
    for i:=0; i<DatabaseConnRetries; i++ {
        db = pg.Connect(&pg.Options{
            Addr: DatabaseAddr,
            Database: DatabaseDatabase,
            User: DatabaseUser,
            Password: DatabasePassword,
        })

        db_con_good := isDbConGood()

        if db_con_good {
            break
        }

        db = nil
        time.Sleep(time.Second)
    }

    if db == nil {
        panic("Can't connect with database")
    }

    // Create schema and default data
    createSchema()
    createDefaultData()

    log.Println("Database connection established")
}

func CloseDB() {
    db.Close()
    log.Println("Database connection closed")
}

func GetDBCon() *pg.DB {
    return db
}

func isDbConGood() bool {
    var n int
    _, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
    return err == nil
}

func createSchema() {
    for _, model := range []interface{}{
        (*User)(nil),
        (*Topic)(nil),
        (*Message)(nil),
        (*Invite)(nil),
    } {
        err := db.CreateTable(model, &orm.CreateTableOptions{
            Temp: DatabaseTemp,
        })

        if err != nil {
            panic("Can't create database schema")
        }
    }
}

func createDefaultData() {
    var err error

    // Users
    user_admin := NewUser("admin", Administrator, "pl")
    err = db.Insert(user_admin)

    if err != nil {
        RiftForumPanic("Can't create admin", err)
    }

    user_basic := NewUser("petermlm", Moderator, "pl")
    err = db.Insert(user_basic)

    if err != nil {
        RiftForumPanic("Can't create petermlm", err)
    }

    rift_bot := NewUser("RiftBot", Basic, "pl")
    err = db.Insert(rift_bot)

    if err != nil {
        RiftForumPanic("Can't create RiftBot", err)
    }
}
