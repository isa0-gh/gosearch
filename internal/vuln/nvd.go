package vuln

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type CVE struct {
	ID          string
	Description string
	Published   string
	Severity    string
	Score       float64
	URL         string
}

const nvdURL = "https://services.nvd.nist.gov/rest/json/cves/2.0"

// NVDSearch searches the NVD CVE database via the NVD REST API v2.
func NVDSearch(query string, pages int) ([]CVE, error) {
	const perPage = 20
	client := &http.Client{}
	var all []CVE

	for p := 0; p < pages; p++ {
		u := fmt.Sprintf("%s?keywordSearch=%s&resultsPerPage=%d&startIndex=%d",
			nvdURL, url.QueryEscape(query), perPage, p*perPage)
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return all, err
		}
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return all, fmt.Errorf("page %d: %w", p+1, err)
		}

		var result struct {
			Vulnerabilities []struct {
				CVE struct {
					ID           string `json:"id"`
					Published    string `json:"published"`
					Descriptions []struct {
						Lang  string `json:"lang"`
						Value string `json:"value"`
					} `json:"descriptions"`
					Metrics struct {
						V31 []struct {
							CVSSData struct {
								BaseScore    float64 `json:"baseScore"`
								BaseSeverity string  `json:"baseSeverity"`
							} `json:"cvssData"`
						} `json:"cvssMetricV31"`
						V30 []struct {
							CVSSData struct {
								BaseScore    float64 `json:"baseScore"`
								BaseSeverity string  `json:"baseSeverity"`
							} `json:"cvssData"`
						} `json:"cvssMetricV30"`
						V2 []struct {
							CVSSData struct {
								BaseScore float64 `json:"baseScore"`
							} `json:"cvssData"`
							BaseSeverity string `json:"baseSeverity"`
						} `json:"cvssMetricV2"`
					} `json:"metrics"`
				} `json:"cve"`
			} `json:"vulnerabilities"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return all, fmt.Errorf("page %d: %w", p+1, err)
		}
		resp.Body.Close()

		if len(result.Vulnerabilities) == 0 {
			break
		}

		for _, item := range result.Vulnerabilities {
			c := item.CVE
			desc := ""
			for _, d := range c.Descriptions {
				if d.Lang == "en" {
					desc = d.Value
					break
				}
			}

			var score float64
			var severity string
			if len(c.Metrics.V31) > 0 {
				score = c.Metrics.V31[0].CVSSData.BaseScore
				severity = c.Metrics.V31[0].CVSSData.BaseSeverity
			} else if len(c.Metrics.V30) > 0 {
				score = c.Metrics.V30[0].CVSSData.BaseScore
				severity = c.Metrics.V30[0].CVSSData.BaseSeverity
			} else if len(c.Metrics.V2) > 0 {
				score = c.Metrics.V2[0].CVSSData.BaseScore
				severity = c.Metrics.V2[0].BaseSeverity
			}

			all = append(all, CVE{
				ID:          c.ID,
				Description: desc,
				Published:   c.Published,
				Severity:    severity,
				Score:       score,
				URL:         fmt.Sprintf("https://nvd.nist.gov/vuln/detail/%s", c.ID),
			})
		}
	}

	return all, nil
}
