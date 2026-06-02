package scrapers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const baseBingURL = "https://www.bing.com/search"

func newBingRequest(query string, first int) (*http.Request, error) {
	url := fmt.Sprintf("%s?q=%s&first=%d&FORM=PERE", baseBingURL, query, first)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	return req, nil
}

// decodeBingURL extracts the real URL from Bing's redirect wrapper.
// The actual URL is base64-encoded in the "u=a1<base64>" query parameter.
func decodeBingURL(raw string) string {
	idx := strings.Index(raw, "&u=a1")
	if idx == -1 {
		return raw
	}
	encoded := raw[idx+5:]
	if end := strings.Index(encoded, "&"); end != -1 {
		encoded = encoded[:end]
	}
	// Bing uses URL-safe base64 without padding
	encoded = strings.ReplaceAll(encoded, "-", "+")
	encoded = strings.ReplaceAll(encoded, "_", "/")
	b, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil {
		return raw
	}
	return string(b)
}

func parseBingResults(doc *goquery.Document) []Result {
	var results []Result
	doc.Find("li.b_algo").Each(func(_ int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h2 a").First().Text())
		href, _ := s.Find("h2 a").First().Attr("href")
		snippet := strings.TrimSpace(s.Find("p").First().Text())
		if title != "" && href != "" {
			results = append(results, Result{
				Title:   title,
				URL:     decodeBingURL(href),
				Snippet: snippet,
			})
		}
	})
	return results
}

// BingSearch fetches up to `pages` pages of results for query.
func BingSearch(query string, pages int) ([]Result, error) {
	client := &http.Client{}
	var allResults []Result

	for p := 0; p < pages; p++ {
		first := p*10 + 1
		req, err := newBingRequest(query, first)
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
		results := parseBingResults(doc)
		if len(results) == 0 {
			break
		}
		allResults = append(allResults, results...)
	}

	return allResults, nil
}
