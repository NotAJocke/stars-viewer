package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/NotAJocke/stars-viewer/internal/routes"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	container := os.Getenv("CONTAINER")

	if container == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Couldn't read env file")
		}
	} else {
		if os.Getenv("GH_TOKEN") == "" {
			log.Fatalln("Please set 'GH_TOKEN' env variable.")
		}
	}

	db, err := sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	app := &routes.App{Db: db}
	router := app.NewRouter()

	port := 8080
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Listening at http://localhost%s\n", addr)

	err = http.ListenAndServe(addr, router)
	if err != nil {
		panic(err)
	}
}
