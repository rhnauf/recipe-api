package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func NewDatabase() (*sql.DB, func() error) {
	port := os.Getenv("DB_PORT")
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		port,
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("SSL_MODE"),
	)

	dbConn, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("ERROR CONNECTING TO DB =>", err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatal("ERROR PING TO DB =>", err)
	}

	log.Println("SUCCESSFULLY CONNECTED TO DB ON PORT =>", port)

	return dbConn, dbConn.Close
}
