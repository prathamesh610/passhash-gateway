package main

import (
	"log"

	"prathameshj.dev/passhash-gateway/db"
	"prathameshj.dev/passhash-gateway/server"
)

func main() {
	db, err := db.NewDataBaseClient()
	if err != nil {
		log.Fatalf("failed to initialize db client: %s", err)
	}

	srv := server.StartServer(db)
	srv.Start()
}
