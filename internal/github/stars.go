package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
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

func Paginate() []StarredRepo {
	token := os.Getenv("GH_TOKEN")
	auth := fmt.Sprintf("Bearer %s", token)

	const nextPattern = `(?i)<([^>]*)>; rel="next"`
	re := regexp.MustCompile(nextPattern)

	client := &http.Client{}

	pagesRemaining := true
	var repos []StarredRepo
	url := "https://api.github.com/user/starred?per_age=100"
	i := 0

	for pagesRemaining {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalln("Couldn't create request to github.")
		}
		req.Header.Set("Authorization", auth)
		req.Header.Set("Accept", "application/vnd.github.v3.star+json")

		res, err := client.Do(req)
		if err != nil {
			log.Fatalln("Error making request:", err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln("Error reading response body:", err)
		}

		parsed := parseStarredRepos(body)
		repos = append(repos, parsed...)

		// linkHeader :=
		linkHeader := res.Header.Get("link")
		log.Printf("link header at i=%d: %s\n", i, linkHeader)
		if linkHeader == "" {
			log.Printf("Link header empty leaving at i=%d\n", i)
			pagesRemaining = false
		}
		if !strings.Contains(linkHeader, `rel="next"`) {
			log.Printf("Link header not empty but no next at i=%d\n", i)
			pagesRemaining = false
		}

		if pagesRemaining {
			match := re.FindStringSubmatch(linkHeader)
			if match == nil {
				log.Fatalf("No match for a next link in '%s'\n", linkHeader)
			}

			url = match[1]
			log.Printf("Link header not empty new url: %s\n", url)
		}

		i += 1
	}

	return repos
}

func parseStarredRepos(body []byte) []StarredRepo {

	var reposRes []starredRepoResponse
	err := json.Unmarshal(body, &reposRes)
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
