package apps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type App struct {
	AppID       string
	Name        string
	Summary     string
	Developer   string
	License     string
	Icon        string
	URL         string
	UpdatedAt   int64
}

const flathubURL = "https://flathub.org/api/v2/search"

// FlathubSearch searches Flathub via the v2 MeiliSearch-backed API.
func FlathubSearch(query string, pages int) ([]App, error) {
	const perPage = 25
	client := &http.Client{}
	var all []App

	for p := 0; p < pages; p++ {
		body, _ := json.Marshal(map[string]any{
			"query":       query,
			"locale":      "en-US",
			"page":        p + 1,
			"hitsPerPage": perPage,
		})
		req, err := http.NewRequest("POST", flathubURL, bytes.NewReader(body))
		if err != nil {
			return all, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return all, fmt.Errorf("page %d: %w", p+1, err)
		}

		var result struct {
			Hits []struct {
				AppID       string `json:"app_id"`
				Name        string `json:"name"`
				Summary     string `json:"summary"`
				DeveloperName string `json:"developer_name"`
				License     string `json:"project_license"`
				Icon        string `json:"icon"`
				UpdatedAt   int64  `json:"updated_at"`
			} `json:"hits"`
			TotalPages int `json:"totalPages"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return all, fmt.Errorf("page %d: %w", p+1, err)
		}
		resp.Body.Close()

		if len(result.Hits) == 0 {
			break
		}

		for _, h := range result.Hits {
			all = append(all, App{
				AppID:     h.AppID,
				Name:      h.Name,
				Summary:   h.Summary,
				Developer: h.DeveloperName,
				License:   h.License,
				Icon:      h.Icon,
				URL:       fmt.Sprintf("https://flathub.org/apps/%s", h.AppID),
				UpdatedAt: h.UpdatedAt,
			})
		}

		if p+1 >= result.TotalPages {
			break
		}
	}

	return all, nil
}
