package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var Instance *sqlx.DB
var err error

func Connect(connectionString string) {
	Instance, err = sqlx.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
		panic("cannot connect to DB")
	}
	log.Println("connceted to DB")
}
