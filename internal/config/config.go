package config

import (
	"fmt"
	"os"
)

func GetDatabaseURL() string {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	environment := os.Getenv("ENVIRONMENT")
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", dbHost, dbUser, dbPassword, dbName)

	if environment == "development" || environment == "docker" {
		return connStr + " sslmode=disable"
	}

	return connStr + " sslmode=verify-full sslrootcert=/var/lib/postgresql/data/server.crt"
}
