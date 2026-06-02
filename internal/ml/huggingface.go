package ml

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// HuggingFaceSearch searches Hugging Face models using its public API.
func HuggingFaceSearch(query string, pages int) ([]Model, error) {
	client := &http.Client{}
	var results []Model

	limit := pages * 25
	escaped := url.QueryEscape(query)
	u := fmt.Sprintf("https://huggingface.co/api/models?search=%s&limit=%d", escaped, limit)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("huggingface API status: %d", resp.StatusCode)
	}

	var apiResponse []struct {
		ID          string   `json:"id"`
		Downloads   int      `json:"downloads"`
		Likes       int      `json:"likes"`
		PipelineTag string   `json:"pipeline_tag"`
		LibraryName string   `json:"library_name"`
		Tags        []string `json:"tags"`
		CreatedAt   string   `json:"createdAt"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	sizeRegex := regexp.MustCompile(`(?i)(?:-|\b)([0-9.]+[Bb])\b`)

	for _, item := range apiResponse {
		href := "https://huggingface.co/" + item.ID

		descParts := []string{}
		if item.PipelineTag != "" {
			descParts = append(descParts, fmt.Sprintf("Pipeline: %s", item.PipelineTag))
		}
		if item.LibraryName != "" {
			descParts = append(descParts, fmt.Sprintf("Library: %s", item.LibraryName))
		}
		descParts = append(descParts, fmt.Sprintf("Likes: %s", formatNumber(item.Likes)))
		description := strings.Join(descParts, " | ")

		var capabilities []string
		if item.PipelineTag != "" {
			capabilities = append(capabilities, item.PipelineTag)
		}
		for _, t := range item.Tags {
			t = strings.ToLower(t)
			if t == "vision" || t == "conversational" || t == "multimodal" || t == "gguf" || t == "safetensors" || t == "uncensored" {
				duplicate := false
				for _, cap := range capabilities {
					if cap == t {
						duplicate = true
						break
					}
				}
				if !duplicate {
					capabilities = append(capabilities, t)
				}
			}
		}

		pulls := formatNumber(item.Downloads)
		tags := fmt.Sprintf("%s Likes", formatNumber(item.Likes))

		size := ""
		matches := sizeRegex.FindStringSubmatch(item.ID)
		if len(matches) > 1 {
			size = strings.ToLower(matches[1])
		}

		updated := ""
		if item.CreatedAt != "" {
			t, err := time.Parse(time.RFC3339Nano, item.CreatedAt)
			if err == nil {
				updated = t.Format("2006-01-02")
			} else {
				if len(item.CreatedAt) >= 10 {
					updated = item.CreatedAt[:10]
				}
			}
		}

		results = append(results, Model{
			Name:         item.ID,
			URL:          href,
			Description:  description,
			Capabilities: capabilities,
			Pulls:        pulls,
			Tags:         tags,
			Size:         size,
			Updated:      updated,
		})
	}

	return results, nil
}

func formatNumber(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000.0)
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000.0)
	}
	return strconv.Itoa(n)
}
