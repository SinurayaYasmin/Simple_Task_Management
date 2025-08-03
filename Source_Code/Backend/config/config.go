package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	name := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || name == "" {
		log.Fatal("Database credentials are not set in .env")
	}

	dsn := fmt.Sprintf("host = %s port=%s user=%s dbname=%s sslmode=disable", host, port, user, name)

	DB, err = sql.Open("postgres", dsn)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("DB unreachable: ", err)
	}

	fmt.Println("Successfully connected to database")

}
