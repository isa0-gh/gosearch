package torrents

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const apibayURL = "https://apibay.org/q.php"

type Torrent struct {
	Name      string
	InfoHash  string
	MagnetURL string
	Seeders   int
	Leechers  int
	Size      int64
	AddedAt   time.Time
	Category  string
	Uploader  string
}

var categories = map[string]string{
	"100": "Audio", "101": "Music", "102": "Audio Books", "103": "Sound Clips", "104": "FLAC", "199": "Audio Other",
	"200": "Video", "201": "Movies", "202": "Movies DVDR", "203": "Music Videos", "204": "Movie Clips",
	"205": "TV Shows", "206": "Handheld", "207": "HD Movies", "208": "HD TV Shows", "209": "3D",
	"210": "CAM/TS", "211": "UHD/4K Movies", "212": "UHD/4K Shows", "299": "Video Other",
	"300": "Applications", "301": "Windows", "302": "Mac", "303": "UNIX", "304": "Handheld",
	"305": "IOS (iPad/iPhone)", "306": "Android", "399": "Applications Other",
	"400": "Games", "401": "PC", "402": "Mac", "403": "PSx", "404": "XBOX360",
	"405": "Wii", "406": "Handheld", "407": "IOS (iPad/iPhone)", "408": "Android", "499": "Games Other",
	"500": "Porn", "501": "Movies", "502": "Movies DVDR", "503": "Pictures", "504": "Games",
	"505": "HD Movies", "506": "Movie Clips", "507": "UHD/4K Movies", "599": "Porn Other",
	"600": "Other", "601": "E-books", "602": "Comics", "603": "Pictures", "604": "Covers",
	"605": "Physibles", "699": "Other Other",
}

var trackers = []string{
	"udp://tracker.opentrackr.org:1337/announce",
	"udp://open.tracker.cl:1337/announce",
	"udp://tracker.openbittorrent.com:6969/announce",
}

func buildMagnet(name, infoHash string) string {
	magnet := fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s", infoHash, url.QueryEscape(name))
	for _, t := range trackers {
		magnet += "&tr=" + url.QueryEscape(t)
	}
	return magnet
}

// PirateBaySearch searches torrents via apibay.org (The Pirate Bay backend API).
// cat=0 means all categories.
func PirateBaySearch(query string, cat int) ([]Torrent, error) {
	u := fmt.Sprintf("%s?q=%s&cat=%d", apibayURL, url.QueryEscape(query), cat)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gosearch")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("apibay: status %d", resp.StatusCode)
	}

	var items []struct {
		Name     string `json:"name"`
		InfoHash string `json:"info_hash"`
		Seeders  string `json:"seeders"`
		Leechers string `json:"leechers"`
		Size     string `json:"size"`
		Added    string `json:"added"`
		Category string `json:"category"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}

	// single "no results" sentinel
	if len(items) == 1 && items[0].Name == "No results returned" {
		return nil, nil
	}

	torrents := make([]Torrent, 0, len(items))
	for _, item := range items {
		seeders, _ := strconv.Atoi(item.Seeders)
		leechers, _ := strconv.Atoi(item.Leechers)
		size, _ := strconv.ParseInt(item.Size, 10, 64)
		added, _ := strconv.ParseInt(item.Added, 10, 64)
		catName := categories[item.Category]
		if catName == "" {
			catName = item.Category
		}
		torrents = append(torrents, Torrent{
			Name:      item.Name,
			InfoHash:  item.InfoHash,
			MagnetURL: buildMagnet(item.Name, item.InfoHash),
			Seeders:   seeders,
			Leechers:  leechers,
			Size:      size,
			AddedAt:   time.Unix(added, 0),
			Category:  catName,
			Uploader:  item.Username,
		})
	}
	return torrents, nil
}
