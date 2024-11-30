package pgDB

import (
	"backend/internal/logger"
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	once     sync.Once
	instance *db
)

type db struct {
	user     string
	password string
	database string
	host     string
	port     string
	db       *sql.DB
}

// Set Connection data for PostgreSQL
func SetConnectionData(user, password, database, host, port string) {
	once.Do(func() {
		instance = &db{
			user:     user,
			password: password,
			database: database,
			host:     host,
			port:     port,
		}
	})
}

// Get instance of PostgreSQL client
func GetDB() *sql.DB {
	if instance == nil {
		logger.GetInstance().Error("database connection is not set", "get PostgreSQL client", nil, nil)
		log.Fatal("PostgreSQL client is nil pointer")
	}

	if instance.db == nil {
		postgresInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			instance.host, instance.port, instance.user, instance.password, instance.database)

		db, err := sql.Open("postgres", postgresInfo)
		if err != nil {
			logger.GetInstance().Error(err.Error(), "connect to PostgreSQL", nil, err)
			log.Fatalf("Could not connect to PostgreSQL: %v", err)
		}

		instance.db = db
	}

	return instance.db
}
