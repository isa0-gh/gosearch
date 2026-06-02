package software

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// SourceForgeSearch searches projects by scraping SourceForge directory.
func SourceForgeSearch(query string, pages int) ([]Repository, error) {
	client := &http.Client{}
	var allRepos []Repository

	for p := 1; p <= pages; p++ {
		u := fmt.Sprintf("https://sourceforge.net/directory/?q=%s&page=%d", url.QueryEscape(query), p)
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return allRepos, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")

		resp, err := client.Do(req)
		if err != nil {
			return allRepos, err
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			return allRepos, err
		}

		var found int
		doc.Find("a[href^='/projects/']").Each(func(_ int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			// only top-level project links: /projects/name/
			parts := strings.Split(strings.Trim(href, "/"), "/")
			if len(parts) != 2 {
				return
			}
			name := s.Find("span[itemprop='name']").Text()
			if name == "" {
				name = parts[1]
			}
			desc := strings.TrimSpace(s.Find("p").Text())
			projectURL := "https://sourceforge.net" + href
			allRepos = append(allRepos, Repository{
				Name:        name,
				URL:         projectURL,
				Description: desc,
			})
			found++
		})

		if found == 0 {
			break
		}
	}

	return allRepos, nil
}
