package config

import (
	"fmt"
	"os"
)

func GetDatabaseURL() string {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	environment := os.Getenv("ENVIRONMENT")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	if environment == "development" || environment == "docker" {
		return connStr + " sslmode=disable"
	}

	return connStr + " sslmode=verify-full"
}
