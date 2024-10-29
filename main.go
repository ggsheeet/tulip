package main

import (
	"log"
	"os"

	"github.com/ggsheet/kerigma/api"
	"github.com/ggsheet/kerigma/app"
	"github.com/ggsheet/kerigma/internal/database"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	isDev := os.Getenv("ENVIRONMENT") == "development"
	if isDev {
		if err := godotenv.Load(".env.local"); err != nil {
			log.Fatal("Error loading .env.local file")
		}
	}

	db, err := database.DBConnection()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	if err := db.Init(); err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}

	e := echo.New()
	e.HideBanner = true

	apiServer := api.NewAPIServer(db)
	apiServer.APIRouter(e)
	log.Println("API routes initialized successfully")

	e.GET("/*", echoWrapHandler(public()))

	app.APPRouter(e)
	log.Println("Application routes initialized successfully")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
