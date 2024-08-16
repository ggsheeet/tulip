package main

import (
	"log"
	"os"

	"github.com/ggsheet/kerigma/api"
	app "github.com/ggsheet/kerigma/app"
	"github.com/ggsheet/kerigma/internal/database"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if os.Getenv("ENVIRONMENT") == "development" {
		if err := godotenv.Load(".env.local"); err != nil {
			log.Fatal("Error loading .env.local file")
		}
	}

	db, err := database.DBConnection()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.HideBanner = true

	apiServer := api.NewAPIServer(db)
	apiServer.APIRouter(e)

	e.GET("/*", echoWrapHandler(public()))

	app.APPRouter(e)

	e.Logger.Fatal(e.Start(":8080"))
}
