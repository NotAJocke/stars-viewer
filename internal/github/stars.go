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

type RepoData struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Url      string `json:"html_url"`
}

type StarredRepo struct {
	Repo      RepoData  `json:"repo"`
	StarredAt time.Time `json:"starred_at"`
}

func GetStarredRepos() []StarredRepo {
	req, err := http.NewRequest("GET", "https://api.github.com/user/starred?per_page=2", nil)
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

	var repos []StarredRepo
	err = json.Unmarshal(body, &repos)
	if err != nil {
		log.Fatalln("Error parsing JSON:", err)
	}

	// for _, repo := range repos {
	// 	log.Printf("Name: %s\n", repo.Repo.Name)
	// 	log.Printf("Full name: %s\n", repo.Repo.FullName)
	// 	log.Printf("URL: %s\n", repo.Repo.Url)
	// 	tz := time.Now().Local().Location()
	// 	log.Printf("Starred at: %s\n", repo.StarredAt.In(tz))
	// }

	return repos
}

func GetStarredReposName() {
}
