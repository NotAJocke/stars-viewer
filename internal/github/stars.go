package github

import (
	"encoding/json"
	"fmt"
	"io"
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

func Paginate() ([]StarredRepo, error) {

	const nextPattern = `(?i)<([^>]*)>; rel="next"`
	re := regexp.MustCompile(nextPattern)

	client := &http.Client{}

	pagesRemaining := true
	var repos []StarredRepo
	url := "https://api.github.com/user/starred?per_page=100"

	for pagesRemaining {
		res, err := requestStarredRepos(client, url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Error reading response body: %w", err)
		}

		parsed, err := parseStarredRepos(body)
		if err != nil {
			return nil, err
		}
		repos = append(repos, parsed...)

		linkHeader := res.Header.Get("link")
		if linkHeader == "" {
			pagesRemaining = false
		}
		if !strings.Contains(linkHeader, `rel="next"`) {
			pagesRemaining = false
		}

		if pagesRemaining {
			match := re.FindStringSubmatch(linkHeader)
			if match == nil {
				return nil, fmt.Errorf("Regex failed finding next link")
			}

			url = match[1]
		}
	}

	return repos, nil
}

func RequestLastStarredRepos(client *http.Client, pages ...int) ([]StarredRepo, error) {
	var page int
	if len(pages) > 0 {
		page = pages[0]
	} else {
		page = 0
	}

	url := fmt.Sprintf("https://api.github.com/user/starred?per_page=30&page=%d", page)

	res, err := requestStarredRepos(client, url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %w", err)
	}

	parsed, err := parseStarredRepos(body)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

func requestStarredRepos(client *http.Client, url string) (*http.Response, error) {
	token := os.Getenv("GH_TOKEN")
	auth := fmt.Sprintf("Bearer %s", token)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Couln't create request to url '%s'", url)
	}
	req.Header.Set("Authorization", auth)
	req.Header.Set("Accept", "application/vnd.github.v3.star+json")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed sending request to url '%s'", url)
	}

	return res, nil
}

func parseStarredRepos(body []byte) ([]StarredRepo, error) {

	var reposRes []starredRepoResponse
	err := json.Unmarshal(body, &reposRes)
	if err != nil {
		return nil, fmt.Errorf("Error parsing json response")
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

	return repos, nil
}

func UnstarRepo(fullName string) error {
	url := fmt.Sprintf("https://api.github.com/user/starred/%s", fullName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("Couln't create request to url '%s'", url)
	}

	token := os.Getenv("GH_TOKEN")
	auth := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed sending request to url '%s'", url)
	}
	defer res.Body.Close()

	return nil
}
