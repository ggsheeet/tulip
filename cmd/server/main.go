package main

import (
	"log"

	"github.com/ggsheet/kerigma/api"
	"github.com/ggsheet/kerigma/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")

	dbInit, err := database.DBConnection()

	if err != nil {
		log.Fatal(err)
	}

	if accErr, bookErr := dbInit.Init(); accErr != nil || bookErr != nil {
		log.Fatal(accErr, bookErr)
	}

	server := api.NewAPIServer(":8080", dbInit, dbInit)
	server.Run()
}

// "html/template"
// "github.com/gorilla/mux"
// e21e8d966b7b4392d35475b146573665c7cdd1b61141a79bd0096e36864ebd02

// var templates *template.Template

// templates = template.Must(template.ParseGlob("src/templates/*.html"))
// r := mux.NewRouter()
// r.HandleFunc("/", indexHandler).Methods("GET")
// http.Handle("/", r)
// http.ListenAndServe(":8080", nil)
// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, "Main")
// }
