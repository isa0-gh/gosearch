package apps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

const (
	algoliaAppID  = "D9HG3G8GS4"
	algoliaAPIKey = "e3369d62b2366b374c54b2c5a2835a00"
	algoliaURL    = "https://D9HG3G8GS4-dsn.algolia.net/1/indexes/brew_all/query"
)

type algoliaResponse struct {
	Hits []struct {
		Hierarchy struct {
			Lvl0 string `json:"lvl0"`
			Lvl1 string `json:"lvl1"`
		} `json:"hierarchy"`
		URL string `json:"url"`
	} `json:"hits"`
}

// HomebrewSearch searches Homebrew formulae and casks using Algolia.
func HomebrewSearch(query string, pages int) ([]App, error) {
	if pages <= 0 {
		pages = 1
	}
	const hitsPerPage = 15
	client := &http.Client{}
	var all []App

	for p := 0; p < pages; p++ {
		payload := map[string]any{
			"params": fmt.Sprintf("query=%s&hitsPerPage=%d&page=%d", query, hitsPerPage, p),
		}
		body, _ := json.Marshal(payload)
		url := fmt.Sprintf("%s?x-algolia-application-id=%s&x-algolia-api-key=%s", algoliaURL, algoliaAppID, algoliaAPIKey)
		
		req, err := http.NewRequest("POST", url, bytes.NewReader(body))
		if err != nil {
			return all, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return all, err
		}

		var algolia algoliaResponse
		if err := json.NewDecoder(resp.Body).Decode(&algolia); err != nil {
			resp.Body.Close()
			return all, err
		}
		resp.Body.Close()

		if len(algolia.Hits) == 0 {
			break
		}

		var wg sync.WaitGroup
		var mu sync.Mutex
		pageResults := make([]App, len(algolia.Hits))
		validIndices := make([]bool, len(algolia.Hits))

		for i, hit := range algolia.Hits {
			if hit.Hierarchy.Lvl0 != "Formulae" && hit.Hierarchy.Lvl0 != "Casks" {
				continue
			}

			wg.Add(1)
			go func(index int, h struct {
				Hierarchy struct {
					Lvl0 string `json:"lvl0"`
					Lvl1 string `json:"lvl1"`
				} `json:"hierarchy"`
				URL string `json:"url"`
			}) {
				defer wg.Done()
				
				// Clean name from lvl1 (e.g. "wget" or "docker-completion (deprecated)")
				name := h.Hierarchy.Lvl1
				if idx := strings.Index(name, " "); idx != -1 {
					name = name[:idx]
				}

				isCask := h.Hierarchy.Lvl0 == "Casks"
				detailURL := ""
				if isCask {
					detailURL = fmt.Sprintf("https://formulae.brew.sh/api/cask/%s.json", name)
				} else {
					detailURL = fmt.Sprintf("https://formulae.brew.sh/api/formula/%s.json", name)
				}

				dResp, err := client.Get(detailURL)
				if err != nil {
					return
				}
				defer dResp.Body.Close()

				if dResp.StatusCode != 200 {
					return
				}

				var app App
				if isCask {
					var cask struct {
						Token    string   `json:"token"`
						Name     []string `json:"name"`
						Desc     string   `json:"desc"`
						Homepage string   `json:"homepage"`
						Version  string   `json:"version"`
					}
					if err := json.NewDecoder(dResp.Body).Decode(&cask); err != nil {
						return
					}
					displayName := cask.Token
					if len(cask.Name) > 0 {
						displayName = cask.Name[0]
					}
					app = App{
						AppID:     cask.Token,
						Name:      displayName,
						Summary:   cask.Desc,
						Developer: "Homebrew Cask",
						License:   "", // Casks don't usually have license in API
						URL:       fmt.Sprintf("https://formulae.brew.sh/cask/%s", cask.Token),
					}
				} else {
					var formula struct {
						Name     string `json:"name"`
						Desc     string `json:"desc"`
						License  string `json:"license"`
						Homepage string `json:"homepage"`
						Versions struct {
							Stable string `json:"stable"`
						} `json:"versions"`
					}
					if err := json.NewDecoder(dResp.Body).Decode(&formula); err != nil {
						return
					}
					app = App{
						AppID:     formula.Name,
						Name:      formula.Name,
						Summary:   formula.Desc,
						Developer: "Homebrew Core",
						License:   formula.License,
						URL:       fmt.Sprintf("https://formulae.brew.sh/formula/%s", formula.Name),
					}
				}

				mu.Lock()
				pageResults[index] = app
				validIndices[index] = true
				mu.Unlock()
			}(i, hit)
		}
		wg.Wait()

		for i, valid := range validIndices {
			if valid {
				all = append(all, pageResults[i])
			}
		}
	}

	return all, nil
}
