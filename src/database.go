package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	_ "github.com/go-pg/pg/orm"
)

var db *pg.DB

func InitDB() {
	// Connect to database, waits one second between attempts
	for i := 0; i < Config.DatabaseConnRetries; i++ {
		db = pg.Connect(&pg.Options{
			Addr:     Config.DatabaseAddr,
			Database: Config.DatabaseDatabase,
			User:     Config.DatabaseUser,
			Password: Config.DatabasePassword,
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

	log.Println("Database connection established")
}

func MigrateCmd(commands []string) {
	oldVersion, newVersion, err := migrations.Run(db, commands...)
	if err != nil {
		panic(err.Error())
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}
}

func MigrationsTableExsits() bool {
	res, _ := db.Exec(`
           SELECT 1
           FROM   information_schema.tables
           WHERE  table_name = 'gopg_migrations';
    `)
	return res.RowsReturned() > 0
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
