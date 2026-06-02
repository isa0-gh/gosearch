package software

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Repository struct {
	Name        string
	URL         string
	Description string
	Stars       int
	Language    string
	UpdatedAt   time.Time
}

// GitHubSearch searches repositories via GitHub API v3.
// Returns up to 100 results (API limit per page * pages).
func GitHubSearch(query string, pages int) ([]Repository, error) {
	const baseURL = "https://api.github.com/search/repositories"
	client := &http.Client{}
	var allRepos []Repository

	for p := 1; p <= pages; p++ {
		u := fmt.Sprintf("%s?q=%s&per_page=100&page=%d", baseURL, url.QueryEscape(query), p)
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return allRepos, err
		}
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("User-Agent", "gosearch")

		resp, err := client.Do(req)
		if err != nil {
			return allRepos, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return allRepos, fmt.Errorf("github api: status %d", resp.StatusCode)
		}

		var result struct {
			Items []struct {
				Name        string    `json:"name"`
				FullName    string    `json:"full_name"`
				HTMLURL     string    `json:"html_url"`
				Description string    `json:"description"`
				Stars       int       `json:"stargazers_count"`
				Language    string    `json:"language"`
				UpdatedAt   time.Time `json:"updated_at"`
			} `json:"items"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return allRepos, err
		}

		if len(result.Items) == 0 {
			break
		}

		for _, item := range result.Items {
			allRepos = append(allRepos, Repository{
				Name:        item.FullName,
				URL:         item.HTMLURL,
				Description: item.Description,
				Stars:       item.Stars,
				Language:    item.Language,
				UpdatedAt:   item.UpdatedAt,
			})
		}
	}

	return allRepos, nil
}
