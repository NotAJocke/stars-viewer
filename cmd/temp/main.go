package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/NotAJocke/stars-viewer/internal/database"
	// "github.com/NotAJocke/stars-viewer/internal/github"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't read env file")
	}

	db, err := sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// repos := github.GetStarredRepos()
	// database.AppendRepos(db, repos)

	repos := database.FetchRepos(db)
	fmt.Println(len(repos))
}
