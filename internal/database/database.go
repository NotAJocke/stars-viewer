package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/NotAJocke/stars-viewer/internal/github"
	_ "github.com/mattn/go-sqlite3"
)

func AppendRepos(db *sql.DB, repos []github.StarredRepo) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	q := `
  INSERT INTO stars (name, full_name, url, starred_at) VALUES 
  `
	var params []interface{}
	for i, repo := range repos {
		if i > 0 {
			q += ","
		}
		q += "(?, ?, ?, ?)"
		params = append(params, repo.Repo.Name, repo.Repo.FullName, repo.Repo.Url, repo.StarredAt.Format(time.RFC3339))

	}

	stmt, err := tx.Prepare(q)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		log.Println("Couldn't insert in DB")
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
}

func FetchRepos(db *sql.DB) []github.StarredRepo {

	q := `SELECT name, full_name, url, starred_at FROM stars`

	rows, err := db.Query(q)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	var repos []github.StarredRepo

	for rows.Next() {
		var repo github.StarredRepo
		var repoData github.RepoData
		var starred_at string

		if err := rows.Scan(&repoData.Name, &repoData.FullName, &repoData.Url, &starred_at); err != nil {
			log.Fatal(err)
		}

		repo.StarredAt, err = time.Parse(time.RFC3339, starred_at)
		if err != nil {
			log.Fatalln(err)
		}
		repo.Repo = repoData

		repos = append(repos, repo)
	}

	if err := rows.Err(); err != nil {
		log.Fatalln(err)
	}

	return repos
}
