package config

import (
	"fmt"
	"os"
)

// GetDatabaseURL returns the database URL based on the environment
func GetDatabaseURL() string {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	environment := os.Getenv("ENVIRONMENT")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	if environment == "development" || environment == "docker" {
		// Apply development-specific configurations if needed
		return connStr + " sslmode=disable"
	}

	return connStr + " sslmode=verify-full"
}

// For the connection between the two containers to happen without errors, I ned to in some way open the database in one of the processes inside the file that contains func main()?
