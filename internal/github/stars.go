package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type repoData struct {
	Name        string   `json:"name"`
	FullName    string   `json:"full_name"`
	Url         string   `json:"html_url"`
	Description string   `json:"description"`
	Language    string   `json:"language"`
	Topics      []string `json:"topics"`
}

type starredRepoResponse struct {
	Repo      repoData  `json:"repo"`
	StarredAt time.Time `json:"starred_at"`
}

type StarredRepo struct {
	Name        string
	FullName    string
	Description string
	Url         string
	StarredAt   time.Time
	Language    string
	Topics      []string
}

func GetStarredRepos() []StarredRepo {
	req, err := http.NewRequest("GET", "https://api.github.com/user/starred?per_page=1", nil)
	if err != nil {
		log.Fatalln("Couldn't create request to github.")
	}

	token := os.Getenv("GH_TOKEN")
	auth := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Accept", "application/vnd.github.v3.star+json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error making request:", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln("Error reading response body:", err)
	}

	var reposRes []starredRepoResponse
	err = json.Unmarshal(body, &reposRes)
	if err != nil {
		log.Fatalln("Error parsing JSON:", err)
	}

	var repos []StarredRepo
	for _, repo := range reposRes {
		repos = append(repos, StarredRepo{
			Name:        repo.Repo.Name,
			Description: repo.Repo.Description,
			FullName:    repo.Repo.FullName,
			Url:         repo.Repo.Url,
			StarredAt:   repo.StarredAt,
			Topics:      repo.Repo.Topics,
			Language:    repo.Repo.Language,
		})
	}

	return repos
}

func UnstarRepo(fullName string) {
	url := fmt.Sprintf("https://api.github.com/user/starred/%s", fullName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatalln("Couldn't create request to github.")
	}

	token := os.Getenv("GH_TOKEN")
	auth := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error making request:", err)
	}
	defer res.Body.Close()
}
