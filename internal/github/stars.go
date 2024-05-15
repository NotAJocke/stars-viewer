package github

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func GetStarredRepos() string {
	req, err := http.NewRequest("GET", "https://api.github.com/user/starred", nil)
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln("Error reading response body:", err)
	}

	return string(body)
}

func GetStarredReposName() {
}
