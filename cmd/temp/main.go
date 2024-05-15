package main

import (
	"encoding/json"
	"log"

	"github.com/NotAJocke/stars-viewer/internal/github"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't read env file")
	}

	repos := github.GetStarredRepos()

	var result []map[string]interface{}
	err = json.Unmarshal([]byte(repos), &result)
	if err != nil {
		log.Fatalln("Error parsing JSON:", err)
	}

	log.Println(len(result))
}
