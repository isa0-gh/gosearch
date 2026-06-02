package scrapers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const baseBraveURL = "https://search.brave.com/search"

func newBraveRequest(query string, offset int) (*http.Request, error) {
	u := fmt.Sprintf("%s?q=%s&offset=%d&spellcheck=0", baseBraveURL, url.QueryEscape(query), offset)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	return req, nil
}

func parseBraveResults(doc *goquery.Document) []Result {
	var results []Result
	doc.Find(`div.snippet[data-type="web"]`).Each(func(_ int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("div.title").First().Text())
		href, _ := s.Find("a").First().Attr("href")
		snippet := strings.TrimSpace(s.Find("div.generic-snippet div.content").First().Text())
		if title != "" && href != "" {
			results = append(results, Result{Title: title, URL: href, Snippet: snippet})
		}
	})
	return results
}

// BraveSearch fetches up to `pages` pages of results for query.
// Each page has ~20 results; pagination uses offset=0, offset=1, offset=2...
func BraveSearch(query string, pages int) ([]Result, error) {
	client := &http.Client{}
	var allResults []Result

	for p := 0; p < pages; p++ {
		req, err := newBraveRequest(query, p)
		if err != nil {
			return allResults, fmt.Errorf("page %d: %w", p+1, err)
		}
		resp, err := client.Do(req)
		if err != nil {
			return allResults, fmt.Errorf("page %d: %w", p+1, err)
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			return allResults, fmt.Errorf("page %d: %w", p+1, err)
		}
		results := parseBraveResults(doc)
		if len(results) == 0 {
			break
		}
		allResults = append(allResults, results...)
	}

	return allResults, nil
}
