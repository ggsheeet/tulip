package main

import (
	"log"
	"os"

	"github.com/ggsheet/tulip/api"
	"github.com/ggsheet/tulip/app"
	"github.com/ggsheet/tulip/internal/database"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/resend/resend-go/v2"
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

	mpAccessToken := os.Getenv("MP_ACCESS_TOKEN")
	cfg, err := config.New(mpAccessToken)
	if err != nil {
		log.Fatalf("failed to initialize Mercadopago client: %v", err)
	}

	resendApiKey := os.Getenv("RESEND_API_KEY")
	msg := resend.NewClient(resendApiKey)

	apiServer := api.NewAPIServer(db, cfg, msg)
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
