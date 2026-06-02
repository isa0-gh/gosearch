package ml

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Model struct {
	Name         string   `json:"name"`
	URL          string   `json:"url"`
	Description  string   `json:"description"`
	Capabilities []string `json:"capabilities"`
	Pulls        string   `json:"pulls"`
	Tags         string   `json:"tags"`
	Size         string   `json:"size"`
	Updated      string   `json:"updated"`
}

// OllamaSearch searches Ollama models using HTML scraping with HTMX headers.
func OllamaSearch(query string, pages int) ([]Model, error) {
	client := &http.Client{}
	var results []Model

	for p := 1; p <= pages; p++ {
		escaped := url.QueryEscape(query)
		u := fmt.Sprintf("https://ollama.com/search?q=%s&page=%d", escaped, p)

		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")
		req.Header.Set("HX-Request", "true")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("ollama API status: %d", resp.StatusCode)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		pageItems := 0
		doc.Find("li[x-test-model]").Each(func(_ int, s *goquery.Selection) {
			pageItems++
			name := strings.TrimSpace(s.Find("span[x-test-search-response-title]").Text())

			href, _ := s.Find("a").First().Attr("href")
			if href != "" && !strings.HasPrefix(href, "http") {
				href = "https://ollama.com" + href
			}

			description := strings.TrimSpace(s.Find("a p").First().Text())

			var capabilities []string
			s.Find("span[x-test-capability]").Each(func(_ int, capSel *goquery.Selection) {
				capText := strings.TrimSpace(capSel.Text())
				if capText != "" {
					capabilities = append(capabilities, capText)
				}
			})

			// Add cloud capability if it exists in the flex list
			s.Find("span").Each(func(_ int, spanSel *goquery.Selection) {
				txt := strings.TrimSpace(spanSel.Text())
				if txt == "cloud" {
					capabilities = append(capabilities, "cloud")
				}
			})

			size := strings.TrimSpace(s.Find("span[x-test-size]").Text())
			pulls := strings.TrimSpace(s.Find("span[x-test-pull-count]").Text())
			tags := strings.TrimSpace(s.Find("span[x-test-tag-count]").Text())
			updated := strings.TrimSpace(s.Find("span[x-test-updated]").Text())

			if name != "" {
				results = append(results, Model{
					Name:         name,
					URL:          href,
					Description:  description,
					Capabilities: capabilities,
					Pulls:        pulls,
					Tags:         tags,
					Size:         size,
					Updated:      updated,
				})
			}
		})

		if pageItems == 0 {
			break
		}
	}

	return results, nil
}
