package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/nanorex07/gobank/api"
	db "github.com/nanorex07/gobank/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgres://root:secret@localhost:5432/gobank?sslmode=disable"
	serverAddress = ":8000"
)

func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database.", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Server start failed: ", err)
	}
}
