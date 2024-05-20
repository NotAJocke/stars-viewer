package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
	tz := time.Now().Local().Location()

	fmt.Println("Number of repos:", len(repos))
	for _, r := range repos {
		fmt.Printf("Id: %d\n", r.Id)
		fmt.Printf("Name: %s\n", r.Name)
		fmt.Printf("Full Name: %s\n", r.FullName)
		fmt.Printf("Url: %s\n", r.Url)
		fmt.Printf("Description: %s\n", r.Description)
		fmt.Printf("Starred at: %s\n", r.StarredAt.In(tz).Format("2006/01/02 - 15:04"))
		fmt.Printf("Labels: %v\n\n", r.Labels)
	}

}
