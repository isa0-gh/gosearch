package scrapers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const baseDDGURL = "https://html.duckduckgo.com/html/"

type Result struct {
	Title   string
	URL     string
	Snippet string
}

// nextPageParams holds the hidden form fields needed to fetch the next page.
type nextPageParams struct {
	S          string
	NextParams string
	V          string
	O          string
	DC         string
	API        string
	VQD        string
	KL         string
}

func newRequest(body url.Values) (*http.Request, error) {
	req, err := http.NewRequest("POST", baseDDGURL, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-GB,en;q=0.9")
	req.Header.Set("Origin", "https://html.duckduckgo.com")
	req.Header.Set("Referer", "https://html.duckduckgo.com/")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")
	return req, nil
}

func parseResults(doc *goquery.Document) ([]Result, *nextPageParams) {
	var results []Result

	doc.Find(".result__body").Each(func(_ int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(".result__a").First().Text())
		href, _ := s.Find(".result__a").First().Attr("href")
		snippet := strings.TrimSpace(s.Find(".result__snippet").First().Text())
		if title != "" && href != "" && !strings.Contains(href, "duckduckgo.com/y.js") {
			results = append(results, Result{Title: title, URL: href, Snippet: snippet})
		}
	})

	// Extract next-page hidden inputs from the form containing the "Next" button.
	var next *nextPageParams
	doc.Find("form").Each(func(_ int, form *goquery.Selection) {
		if form.Find("input[type=submit][value=Next]").Length() == 0 {
			return
		}
		next = &nextPageParams{
			S:          form.Find("input[name=s]").AttrOr("value", ""),
			NextParams: form.Find("input[name=nextParams]").AttrOr("value", ""),
			V:          form.Find("input[name=v]").AttrOr("value", ""),
			O:          form.Find("input[name=o]").AttrOr("value", ""),
			DC:         form.Find("input[name=dc]").AttrOr("value", ""),
			API:        form.Find("input[name=api]").AttrOr("value", ""),
			VQD:        form.Find("input[name=vqd]").AttrOr("value", ""),
			KL:         form.Find("input[name=kl]").AttrOr("value", ""),
		}
	})

	return results, next
}

// DuckDuckGoSearch fetches up to `pages` pages of results for query.
func DuckDuckGoSearch(query string, pages int) ([]Result, error) {
	client := &http.Client{}
	var allResults []Result

	// First page
	body := url.Values{"q": {query}, "b": {""}}
	req, err := newRequest(body)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	results, next := parseResults(doc)
	allResults = append(allResults, results...)

	// Subsequent pages
	for p := 1; p < pages && next != nil; p++ {
		body = url.Values{
			"q":          {query},
			"s":          {next.S},
			"nextParams": {next.NextParams},
			"v":          {next.V},
			"o":          {next.O},
			"dc":         {next.DC},
			"api":        {next.API},
			"vqd":        {next.VQD},
			"kl":         {next.KL},
		}
		req, err = newRequest(body)
		if err != nil {
			return allResults, fmt.Errorf("page %d: %w", p+1, err)
		}
		resp, err = client.Do(req)
		if err != nil {
			return allResults, fmt.Errorf("page %d: %w", p+1, err)
		}
		doc, err = goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			return allResults, fmt.Errorf("page %d: %w", p+1, err)
		}
		results, next = parseResults(doc)
		allResults = append(allResults, results...)
	}

	return allResults, nil
}
