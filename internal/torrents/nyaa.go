package torrents

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type NyaaTorrent struct {
	Name      string
	URL       string
	MagnetURL string
	Size      string
	Category  string
	Seeders   int
	Leechers  int
	Downloads int
	AddedAt   time.Time
}

// NyaaSearch searches nyaa.si. f=filter(0=no filter), c=category(0_0=all), p=page.
func NyaaSearch(query string, pages int) ([]NyaaTorrent, error) {
	client := &http.Client{}
	var all []NyaaTorrent

	for p := 1; p <= pages; p++ {
		u := fmt.Sprintf("https://nyaa.si/?f=0&c=0_0&q=%s&p=%d", url.QueryEscape(query), p)
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return all, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")

		resp, err := client.Do(req)
		if err != nil {
			return all, fmt.Errorf("page %d: %w", p, err)
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			return all, fmt.Errorf("page %d: %w", p, err)
		}

		var found int
		doc.Find("table.torrent-list tbody tr").Each(func(_ int, s *goquery.Selection) {
			tds := s.Find("td")
			if tds.Length() < 8 {
				return
			}

			category := tds.Eq(0).Find("a").AttrOr("title", "")
			nameEl := tds.Eq(1).Find("a").Last()
			name := strings.TrimSpace(nameEl.Text())
			href, _ := nameEl.Attr("href")
			magnet, _ := tds.Eq(2).Find("a[href^='magnet:']").Attr("href")
			size := strings.TrimSpace(tds.Eq(3).Text())
			ts, _ := strconv.ParseInt(tds.Eq(4).AttrOr("data-timestamp", "0"), 10, 64)
			seeders, _ := strconv.Atoi(strings.TrimSpace(tds.Eq(5).Text()))
			leechers, _ := strconv.Atoi(strings.TrimSpace(tds.Eq(6).Text()))
			downloads, _ := strconv.Atoi(strings.TrimSpace(tds.Eq(7).Text()))

			if name == "" {
				return
			}
			all = append(all, NyaaTorrent{
				Name:      name,
				URL:       "https://nyaa.si" + href,
				MagnetURL: magnet,
				Size:      size,
				Category:  category,
				Seeders:   seeders,
				Leechers:  leechers,
				Downloads: downloads,
				AddedAt:   time.Unix(ts, 0),
			})
			found++
		})

		if found == 0 {
			break
		}
	}

	return all, nil
}
