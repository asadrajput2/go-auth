package main

import (
	"log"

	"github.com/asadrajput2/go-auth/pkg/http/rest"
	"github.com/asadrajput2/go-auth/pkg/postgres"
)

func main() {
	db, err := postgres.Connect()

	if err != nil {
		log.Fatal(err)
	}

	rest.ReqHandler(db)
}
