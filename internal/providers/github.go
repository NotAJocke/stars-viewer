package providers

import (
	"database/sql"
	"net/http"

	"github.com/NotAJocke/stars-viewer/internal/database"
	"github.com/NotAJocke/stars-viewer/internal/github"
)

func UpdateCache(db *sql.DB) ([]github.StarredRepo, error) {
	client := &http.Client{}
	lastStarred := database.GetLastStarred(db)

	var reposToAdd []github.StarredRepo
	match := false
	i := 1
	for match != true {
		repos, err := github.RequestLastStarredRepos(client, i)
		if err != nil {
			return nil, err
		}
		for _, repo := range repos {
			if repo.FullName == lastStarred.FullName {
				match = true
				break
			}

			reposToAdd = append(reposToAdd, repo)
		}

		i += 1
	}

	database.AddRepos(db, reposToAdd)

	return reposToAdd, nil
}

func SearchRepos(db *sql.DB, query string, limit int) {
	database.SearchRepos(db, query, limit)
}

func GetStarredRepos(db *sql.DB, limit int) []github.StarredRepo {
	return database.GetRepos(db, limit)
}

func UnstarRepo(db *sql.DB, fullName string) {
	github.UnstarRepo(fullName)

	database.DeleteRepo(db, fullName)
}
