package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/nanorex07/gobank/api"
	db "github.com/nanorex07/gobank/db/sqlc"
	"github.com/nanorex07/gobank/util"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load env variables.", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database.", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Server start failed: ", err)
	}
}
