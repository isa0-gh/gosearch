package academic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

const openAlexURL = "https://api.openalex.org/works"

// reconstructAbstract rebuilds plain text from OpenAlex inverted index format.
func reconstructAbstract(inv map[string][]int) string {
	words := make(map[int]string)
	for word, positions := range inv {
		for _, pos := range positions {
			words[pos] = word
		}
	}
	keys := make([]int, 0, len(words))
	for k := range words {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = words[k]
	}
	return strings.Join(parts, " ")
}

// OpenAlexSearch searches OpenAlex for academic works.
func OpenAlexSearch(query string, pages int) ([]Paper, error) {
	const perPage = 25
	client := &http.Client{}
	var all []Paper

	for p := 1; p <= pages; p++ {
		u := fmt.Sprintf("%s?search=%s&per_page=%d&page=%d&mailto=gosearch",
			openAlexURL, url.QueryEscape(query), perPage, p)
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return all, err
		}
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return all, fmt.Errorf("page %d: %w", p, err)
		}

		var result struct {
			Results []struct {
				Title           string          `json:"title"`
				DOI             string          `json:"doi"`
				PrimaryLocation struct {
					LandingPageURL string `json:"landing_page_url"`
				} `json:"primary_location"`
				Authorships []struct {
					Author struct {
						DisplayName string `json:"display_name"`
					} `json:"author"`
				} `json:"authorships"`
				AbstractInvertedIndex map[string][]int `json:"abstract_inverted_index"`
				Type                  string           `json:"type"`
			} `json:"results"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return all, fmt.Errorf("page %d: %w", p, err)
		}
		resp.Body.Close()

		if len(result.Results) == 0 {
			break
		}

		for _, r := range result.Results {
			link := r.PrimaryLocation.LandingPageURL
			if link == "" {
				link = r.DOI
			}
			var names []string
			for _, a := range r.Authorships {
				if n := a.Author.DisplayName; n != "" {
					names = append(names, n)
				}
			}
			all = append(all, Paper{
				Title:    r.Title,
				URL:      link,
				Authors:  strings.Join(names, ", "),
				Abstract: reconstructAbstract(r.AbstractInvertedIndex),
				Type:     r.Type,
			})
		}
	}

	return all, nil
}
