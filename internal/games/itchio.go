package games

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ItchGame struct {
	GameID       string
	Title        string
	URL          string
	Author       string
	AuthorURL    string
	Description  string
	Genre        string
	ThumbnailURL string
	Platforms    ItchPlatforms
	Rating       *ItchRating
}

type ItchPlatforms struct {
	Windows bool
	MacOS   bool
	Linux   bool
	Web     bool
	Android bool
}

type ItchRating struct {
	Average        float64
	Total          int
	StarPercentage float64
}

var (
	ratingRe = regexp.MustCompile(`([\d.]+).*?from\s+(\d+)`)
	widthRe  = regexp.MustCompile(`width:\s*([\d.]+)%`)
)

func ItchSearch(query string, pages int) ([]ItchGame, error) {
	client := &http.Client{}
	var results []ItchGame

	for p := 1; p <= pages; p++ {
		u := fmt.Sprintf("https://itch.io/search?q=%s&page=%d", url.QueryEscape(query), p)
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return results, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")

		resp, err := client.Do(req)
		if err != nil {
			return results, fmt.Errorf("page %d: %w", p, err)
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return results, fmt.Errorf("itch.io status: %d", resp.StatusCode)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			return results, fmt.Errorf("page %d: %w", p, err)
		}

		count := 0
		doc.Find("div.game_cell").Each(func(_ int, s *goquery.Selection) {
			g := ItchGame{}
			g.GameID, _ = s.Attr("data-game_id")

			if a := s.Find(".game_title .title"); a.Length() > 0 {
				g.Title = strings.TrimSpace(a.Text())
				g.URL, _ = a.Attr("href")
				if !strings.HasPrefix(g.URL, "http") {
					g.URL = "https://itch.io" + g.URL
				}
			}
			if a := s.Find(".game_author a"); a.Length() > 0 {
				g.Author = strings.TrimSpace(a.Text())
				g.AuthorURL, _ = a.Attr("href")
			}
			if d := s.Find(".game_text"); d.Length() > 0 {
				g.Description = d.AttrOr("title", strings.TrimSpace(d.Text()))
			}
			g.Genre = strings.TrimSpace(s.Find(".game_genre").Text())

			if img := s.Find(".game_thumb img"); img.Length() > 0 {
				g.ThumbnailURL = img.AttrOr("data-lazy_src", img.AttrOr("src", ""))
			}

			pt := s.Find(".game_platform").Text()
			g.Platforms = ItchPlatforms{
				Windows: strings.Contains(pt, "icon-windows8") || s.Find(".icon-windows8").Length() > 0,
				MacOS:   strings.Contains(pt, "icon-apple") || s.Find(".icon-apple").Length() > 0,
				Linux:   strings.Contains(pt, "icon-linux") || s.Find(".icon-linux").Length() > 0,
				Web:     strings.Contains(pt, "icon-html5") || s.Find(".icon-html5").Length() > 0,
				Android: strings.Contains(pt, "icon-android") || s.Find(".icon-android").Length() > 0,
			}
			// platform icons are spans, check via html
			s.Find(".game_platform").Each(func(_ int, ps *goquery.Selection) {
				html, _ := ps.Html()
				g.Platforms.Windows = strings.Contains(html, "icon-windows8")
				g.Platforms.MacOS = strings.Contains(html, "icon-apple")
				g.Platforms.Linux = strings.Contains(html, "icon-linux")
				g.Platforms.Web = strings.Contains(html, "icon-html5")
				g.Platforms.Android = strings.Contains(html, "icon-android")
			})

			if rd := s.Find(".game_rating"); rd.Length() > 0 {
				tooltip := rd.AttrOr("data-tooltip", "")
				if m := ratingRe.FindStringSubmatch(tooltip); m != nil {
					avg, _ := strconv.ParseFloat(m[1], 64)
					total, _ := strconv.Atoi(m[2])
					var pct float64
					if sf := rd.Find(".star_fill"); sf.Length() > 0 {
						if wm := widthRe.FindStringSubmatch(sf.AttrOr("style", "")); wm != nil {
							pct, _ = strconv.ParseFloat(wm[1], 64)
						}
					}
					g.Rating = &ItchRating{Average: avg, Total: total, StarPercentage: pct}
				}
			}

			if g.Title != "" {
				results = append(results, g)
				count++
			}
		})

		if count == 0 {
			break
		}
	}

	return results, nil
}
