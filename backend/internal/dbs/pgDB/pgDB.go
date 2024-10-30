package pgDB

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

var db *sql.DB

func Connect(user, password, database, host, port string) {
	postgresInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, database)

	var err error
	db, err = sql.Open("postgres", postgresInfo)
	if err != nil {
		log.Fatalf("Could not connect to PostgreSQL: %v", err)
	}
}

func GetDB() *sql.DB {
	return db
}