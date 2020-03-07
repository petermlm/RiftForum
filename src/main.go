package main

import (
	"flag"
	"log"
	"net/http"
)

func process_migrations() {
	if !MigrationsTableExsits() {
		MigrateCmd([]string{"init"})
	}
	MigrateCmd([]string{})
}

func riftforum() {
	log.Println("Rift Forum Starting")

	InitTmpl()
	InitSers()
	InitAuth()
	InitBots()

	log.Println("Starting Server")
	router := CreateRouter()
	status := http.ListenAndServe(HostAndPort, router)
	log.Fatal(status)
}

func main() {
	migrate := flag.Bool("migrate", false, "Run migrations")
	flag.Parse()

	InitDB()
	defer CloseDB()

	if *migrate {
		process_migrations()
	} else {
		riftforum()
	}
}
