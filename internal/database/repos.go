package database

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/NotAJocke/stars-viewer/internal/github"
	_ "github.com/mattn/go-sqlite3"
)

type StarredRepo struct {
	Id          int
	Name        string
	FullName    string
	Url         string
	StarredAt   time.Time
	Description string
	Labels      []string
}

func AppendRepos(db *sql.DB, repos []github.StarredRepo) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	q := `
  INSERT INTO stars (name, full_name, url, starred_at, description) VALUES 
  `
	var params []interface{}
	for i, repo := range repos {
		if i > 0 {
			q += ","
		}
		q += "(?, ?, ?, ?, ?)"
		params = append(params, repo.Name, repo.FullName, repo.Url, repo.StarredAt.Format(time.RFC3339), repo.Description)
	}

	stmt, err := tx.Prepare(q)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
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

func FetchRepos(db *sql.DB) []StarredRepo {

	q := `SELECT s.id, s.name, s.full_name, s.url, s.starred_at, s.description, IFNULL(GROUP_CONCAT(l.name, ', '), '') AS labels FROM stars s
  LEFT JOIN stars_labels sl ON s.id = sl.star_id
  LEFT JOIN labels l ON sl.label_id = l.id
  GROUP BY s.id, s.full_name;
  `

	rows, err := db.Query(q)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	var repos []StarredRepo

	for rows.Next() {
		var repo StarredRepo
		var starred_at string
		var labels string

		if err := rows.Scan(&repo.Id, &repo.Name, &repo.FullName, &repo.Url, &starred_at, &repo.Description, &labels); err != nil {
			log.Fatal(err)
		}

		repo.StarredAt, err = time.Parse(time.RFC3339, starred_at)
		if err != nil {
			log.Fatalln(err)
		}

		if labels == "" {
			repo.Labels = []string{}
		} else {
			repo.Labels = strings.Split(labels, ", ")
		}

		repos = append(repos, repo)
	}

	if err := rows.Err(); err != nil {
		log.Fatalln(err)
	}

	return repos
}

func DeleteRepoById(db *sql.DB, id int) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}

	q1 := `DELETE FROM stars_labels WHERE star_id=?;`
	q2 := `DELETE FROM stars WHERE id=?;`

	stmt, err := tx.Prepare(q1)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	stmt, err = tx.Prepare(q2)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
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
