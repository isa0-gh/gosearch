package academic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Paper struct {
	Title    string
	URL      string
	Authors  string
	Abstract string
	Type     string
}

const nasaBaseURL = "https://ntrs.nasa.gov/api/citations/search"

// NASASearch searches the NASA Technical Reports Server (NTRS) API.
func NASASearch(query string, pages int) ([]Paper, error) {
	const rows = 25
	client := &http.Client{}
	var all []Paper

	for p := 0; p < pages; p++ {
		u := fmt.Sprintf("%s?q=%s&rows=%d&page=%d", nasaBaseURL, url.QueryEscape(query), rows, p)
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return all, err
		}
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return all, fmt.Errorf("page %d: %w", p+1, err)
		}

		var result struct {
			Results []struct {
				ID                 int    `json:"id"`
				Title              string `json:"title"`
				Abstract           string `json:"abstract"`
				StiTypeDetails     string `json:"stiTypeDetails"`
				AuthorAffiliations []struct {
					Meta struct {
						Author struct {
							Name string `json:"name"`
						} `json:"author"`
					} `json:"meta"`
				} `json:"authorAffiliations"`
			} `json:"results"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return all, fmt.Errorf("page %d: %w", p+1, err)
		}
		resp.Body.Close()

		if len(result.Results) == 0 {
			break
		}

		for _, r := range result.Results {
			var names []string
			for _, a := range r.AuthorAffiliations {
				if n := a.Meta.Author.Name; n != "" {
					names = append(names, n)
				}
			}
			all = append(all, Paper{
				Title:    r.Title,
				URL:      fmt.Sprintf("https://ntrs.nasa.gov/citations/%d", r.ID),
				Authors:  strings.Join(names, ", "),
				Abstract: r.Abstract,
				Type:     r.StiTypeDetails,
			})
		}
	}

	return all, nil
}
