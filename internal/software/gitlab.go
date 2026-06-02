package software

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// GitLabSearch searches projects via GitLab API v4.
func GitLabSearch(query string, pages int) ([]Repository, error) {
	const baseURL = "https://gitlab.com/api/v4/projects"
	client := &http.Client{}
	var allRepos []Repository

	for p := 1; p <= pages; p++ {
		u := fmt.Sprintf("%s?search=%s&per_page=100&page=%d&order_by=last_activity_at", baseURL, url.QueryEscape(query), p)
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return allRepos, err
		}
		req.Header.Set("User-Agent", "gosearch")

		resp, err := client.Do(req)
		if err != nil {
			return allRepos, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return allRepos, fmt.Errorf("gitlab api: status %d", resp.StatusCode)
		}

		var items []struct {
			PathWithNamespace string    `json:"path_with_namespace"`
			WebURL            string    `json:"web_url"`
			Description       string    `json:"description"`
			StarCount         int       `json:"star_count"`
			Language          string    `json:"-"` // not in list endpoint
			LastActivityAt    time.Time `json:"last_activity_at"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
			return allRepos, err
		}

		if len(items) == 0 {
			break
		}

		for _, item := range items {
			allRepos = append(allRepos, Repository{
				Name:        item.PathWithNamespace,
				URL:         item.WebURL,
				Description: item.Description,
				Stars:       item.StarCount,
				UpdatedAt:   item.LastActivityAt,
			})
		}
	}

	return allRepos, nil
}
