package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/NotAJocke/stars-viewer/internal/github"
	_ "github.com/mattn/go-sqlite3"
)

var (
	QueryFailedError error = fmt.Errorf("Query failed")
)

func GetLastStarred(db *sql.DB) github.StarredRepo {
	q := `SELECT name, full_name, description, url, starred_at, language, topics 
  FROM stars ORDER BY starred_at DESC LIMIT 1;`

	rows, err := db.Query(q)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	var repo github.StarredRepo
	var starred_at string
	var topics string

	if rows.Next() {
		if err := rows.Scan(
			&repo.Name,
			&repo.FullName,
			&repo.Description,
			&repo.Url,
			&starred_at,
			&repo.Language,
			&topics,
		); err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Fatalln("No rows returned")

	}

	repo.StarredAt, err = time.Parse(time.RFC3339, starred_at)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal([]byte(topics), &repo.Topics)
	if err != nil {
		log.Fatalln(err)
	}

	return repo
}

func GetRepos(db *sql.DB, limit int) []github.StarredRepo {
	var q string

	if limit != 0 {
		q = `SELECT 
    name, description, full_name, url, starred_at, topics, language 
    FROM stars ORDER BY starred_at DESC LIMIT ?;`
	} else {
		q = `SELECT 
    name, description, full_name, url, starred_at, topics, language 
    FROM stars ORDER BY starred_at DESC;`
	}

	var rows *sql.Rows
	var err error
	if limit != 0 {
		rows, err = db.Query(q, limit)
	} else {
		rows, err = db.Query(q)
	}

	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	var repos []github.StarredRepo

	for rows.Next() {
		var repo github.StarredRepo
		var starred_at string
		var topics string

		if err := rows.Scan(
			&repo.Name,
			&repo.Description,
			&repo.FullName,
			&repo.Url,
			&starred_at,
			&topics,
			&repo.Language,
		); err != nil {
			log.Fatalln(err)
		}

		repo.StarredAt, err = time.Parse(time.RFC3339, starred_at)
		if err != nil {
			log.Fatalln(err)
		}

		err = json.Unmarshal([]byte(topics), &repo.Topics)
		if err != nil {
			log.Fatalln(err)
		}

		repos = append(repos, repo)
	}

	if err := rows.Err(); err != nil {
		log.Fatalln(err)
	}

	return repos
}

func AddRepos(db *sql.DB, repos []github.StarredRepo) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	q := `INSERT INTO stars (name, description, full_name, url, starred_at, topics, language) VALUES (?, ?, ?, ?, ?, ?, ?);`

	stmt, err := tx.Prepare(q)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, repo := range repos {
		topics, err := json.Marshal(repo.Topics)
		if err != nil {
			log.Fatalln(err)
		}

		_, err = stmt.Exec(repo.Name, repo.Description, repo.FullName, repo.Url, repo.StarredAt.Format(time.RFC3339), topics, repo.Language)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}
}

func DeleteRepo(db *sql.DB, repoFullName string) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	q := `DELETE FROM stars WHERE full_name=?`

	stmt, err := tx.Prepare(q)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(repoFullName)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}
}

func SearchRepos(db *sql.DB, query string, limit int) []github.StarredRepo {
	q := `SELECT
    name, description, full_name, url, starred_at, topics, language
    FROM stars WHERE
    name LIKE ? OR 
    full_name LIKE ? OR
    description LIKE ? OR 
    url LIKE ? OR
    starred_at LIKE ? OR
    language LIKE ? OR
    topics LIKE ?
    ORDER BY starred_at DESC
    LIMIT ?;
  `

	rows, err := db.Query(
		q,
		query,
		query,
		query,
		query,
		query,
		query,
		query,
		limit,
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	var repos []github.StarredRepo

	for rows.Next() {
		var repo github.StarredRepo
		var starred_at string
		var topics string

		if err := rows.Scan(
			&repo.Name,
			&repo.Description,
			&repo.FullName,
			&repo.Url,
			&starred_at,
			&topics,
			&repo.Language,
		); err != nil {
			log.Fatalln(err)
		}

		repo.StarredAt, err = time.Parse(time.RFC3339, starred_at)
		if err != nil {
			log.Fatalln(err)
		}

		err = json.Unmarshal([]byte(topics), &repo.Topics)
		if err != nil {
			log.Fatalln(err)
		}

		repos = append(repos, repo)
	}

	if err := rows.Err(); err != nil {
		log.Fatalln(err)
	}

	return repos
}
