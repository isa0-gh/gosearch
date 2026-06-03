package games

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Game struct {
	AppID           string
	Title           string
	URL             string
	ImageURL        string
	ReleaseDate     string
	Price           string
	OriginalPrice   string
	DiscountPercent string
	ReviewSummary   string
	ReviewClass     string
	Platforms       []string
}

func SteamSearch(query string, pages int) ([]Game, error) {
	client := &http.Client{}
	var results []Game

	for p := 1; p <= pages; p++ {
		u := fmt.Sprintf(
			"https://store.steampowered.com/search?term=%s&page=%d",
			url.QueryEscape(query), p,
		)
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return results, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")

		resp, err := client.Do(req)
		if err != nil {
			return results, fmt.Errorf("page %d: %w", p, err)
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return results, fmt.Errorf("steam status: %d", resp.StatusCode)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			return results, fmt.Errorf("page %d: %w", p, err)
		}

		count := 0
		doc.Find("a.search_result_row").Each(func(_ int, s *goquery.Selection) {
			game := Game{}

			game.AppID, _ = s.Attr("data-ds-appid")
			game.URL, _ = s.Attr("href")

			if img := s.Find("img").First(); img.Length() > 0 {
				game.ImageURL, _ = img.Attr("src")
			}
			game.Title = strings.TrimSpace(s.Find("span.title").Text())
			game.ReleaseDate = strings.TrimSpace(s.Find("div.search_released").Text())
			game.Price = strings.TrimSpace(s.Find("div.discount_final_price").Text())
			game.OriginalPrice = strings.TrimSpace(s.Find("div.discount_original_price").Text())
			game.DiscountPercent = strings.TrimSpace(s.Find("div.discount_pct").Text())

			if rev := s.Find("span.search_review_summary"); rev.Length() > 0 {
				game.ReviewSummary, _ = rev.Attr("data-tooltip-html")
				game.ReviewSummary = strings.ReplaceAll(game.ReviewSummary, "<br>", " ")
				classes := strings.Fields(rev.AttrOr("class", ""))
				if len(classes) > 1 {
					game.ReviewClass = classes[1]
				}
			}

			s.Find("div.search_platforms span.platform_img").Each(func(_ int, ps *goquery.Selection) {
				cls := strings.Fields(ps.AttrOr("class", ""))
				if len(cls) > 1 {
					game.Platforms = append(game.Platforms, cls[1])
				}
			})

			if game.Title != "" {
				results = append(results, game)
				count++
			}
		})

		if count == 0 {
			break
		}
	}

	return results, nil
}
