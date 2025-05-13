package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type RepoInfo struct {
	Size int `json:"size"` // Size in KB
}

func getRepoSize(repoURL string) (int, error) {
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return 0, fmt.Errorf("invalid URL: %v", err)
	}

	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid GitHub repo URL format")
	}
	owner, repo := parts[0], parts[1]
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "go-client")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("GitHub API returned status code %d", resp.StatusCode)
	}

	var repoInfo RepoInfo
	err = json.NewDecoder(resp.Body).Decode(&repoInfo)
	if err != nil {
		return 0, err
	}

	return repoInfo.Size, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: github-size <github-repo-url>")
		os.Exit(1)
	}

	repoURL := os.Args[1]
	sizeKB, err := getRepoSize(repoURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Repository size: %d KB (~%.2f MB)\n", sizeKB, float64(sizeKB)/1024)
}
