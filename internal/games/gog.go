package games

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type GogGame struct {
	ID              string
	Slug            string
	Title           string
	URL             string
	ImageURL        string
	ReleaseDate     string
	Price           string
	OriginalPrice   string
	DiscountPercent string
	Developers      []string
	Publishers      []string
	Platforms       []string
	Genres          []string
	Tags            []string
}

func GogSearch(query string, pages int) ([]GogGame, error) {
	client := &http.Client{}
	var results []GogGame

	for p := 1; p <= pages; p++ {
		u := fmt.Sprintf(
			"https://catalog.gog.com/v1/catalog?limit=48&query=like:%s&order=desc:score&productType=in:game,pack,dlc,extras&page=%d&countryCode=TR&locale=en-US",
			url.QueryEscape(query), p,
		)

		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return results, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return results, fmt.Errorf("page %d: %w", p, err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return results, fmt.Errorf("gog API status: %d", resp.StatusCode)
		}

		var apiResponse struct {
			Products []struct {
				ID                string   `json:"id"`
				Slug              string   `json:"slug"`
				Title             string   `json:"title"`
				StoreLink         string   `json:"storeLink"`
				CoverHorizontal   string   `json:"coverHorizontal"`
				CoverVertical     string   `json:"coverVertical"`
				ReleaseDate       string   `json:"releaseDate"`
				Developers        []string `json:"developers"`
				Publishers        []string `json:"publishers"`
				OperatingSystems  []string `json:"operatingSystems"`
				Genres            []struct {
					Name string `json:"name"`
					Slug string `json:"slug"`
				} `json:"genres"`
				Tags []struct {
					Name string `json:"name"`
					Slug string `json:"slug"`
				} `json:"tags"`
				Price *struct {
					Final    string  `json:"final"`
					Base     string  `json:"base"`
					Discount *string `json:"discount"`
				} `json:"price"`
			} `json:"products"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			resp.Body.Close()
			return results, fmt.Errorf("page %d: %w", p, err)
		}
		resp.Body.Close()

		if len(apiResponse.Products) == 0 {
			break
		}

		for _, prod := range apiResponse.Products {
			g := GogGame{
				ID:            prod.ID,
				Slug:          prod.Slug,
				Title:         prod.Title,
				URL:           prod.StoreLink,
				ReleaseDate:   prod.ReleaseDate,
				Developers:    prod.Developers,
				Publishers:    prod.Publishers,
				Platforms:     prod.OperatingSystems,
			}

			if prod.CoverHorizontal != "" {
				g.ImageURL = prod.CoverHorizontal
			} else if prod.CoverVertical != "" {
				g.ImageURL = prod.CoverVertical
			}

			for _, genre := range prod.Genres {
				g.Genres = append(g.Genres, genre.Name)
			}
			for _, tag := range prod.Tags {
				g.Tags = append(g.Tags, tag.Name)
			}

			if prod.Price != nil {
				g.Price = prod.Price.Final
				g.OriginalPrice = prod.Price.Base
				if prod.Price.Discount != nil && *prod.Price.Discount != "" {
					g.DiscountPercent = *prod.Price.Discount
				}
			}

			if g.URL == "" && g.Slug != "" {
				g.URL = "https://www.gog.com/en/game/" + g.Slug
			}

			if g.Title != "" {
				results = append(results, g)
			}
		}
	}

	return results, nil
}
