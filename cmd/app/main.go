package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NotAJocke/stars-viewer/internal/routes"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't read env file")
	}

	router := routes.NewRouter()

	port := 8080
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Listening at http://localhost%s\n", addr)

	err = http.ListenAndServe(addr, router)
	if err != nil {
		panic(err)
	}
}
