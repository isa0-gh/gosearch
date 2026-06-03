package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/isa0-gh/gosearch/internal/academic"
	"github.com/isa0-gh/gosearch/internal/apps"
	"github.com/isa0-gh/gosearch/internal/games"
	"github.com/isa0-gh/gosearch/internal/ml"
	"github.com/isa0-gh/gosearch/internal/scrapers"
	"github.com/isa0-gh/gosearch/internal/software"
	"github.com/isa0-gh/gosearch/internal/torrents"
	"github.com/isa0-gh/gosearch/internal/vuln"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func queryParam(r *http.Request, key string) string { return r.URL.Query().Get(key) }

func pagesParam(r *http.Request) int {
	n, err := strconv.Atoi(r.URL.Query().Get("pages"))
	if err != nil || n < 1 {
		return 1
	}
	return n
}

// GET /api/v1/web?q=...&engine=ddg|bing|brave&pages=1
func handleWeb(w http.ResponseWriter, r *http.Request) {
	q := queryParam(r, "q")
	if q == "" {
		writeError(w, "q is required", http.StatusBadRequest)
		return
	}
	pages := pagesParam(r)

	var results []scrapers.Result
	var err error

	switch queryParam(r, "engine") {
	case "bing":
		results, err = scrapers.BingSearch(q, pages)
	case "brave":
		results, err = scrapers.BraveSearch(q, pages)
	default: // ddg
		results, err = scrapers.DuckDuckGoSearch(q, pages)
	}

	if err != nil {
		writeError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, results)
}

// GET /api/v1/software?q=...&source=github|gitlab|sourceforge&pages=1
func handleSoftware(w http.ResponseWriter, r *http.Request) {
	q := queryParam(r, "q")
	if q == "" {
		writeError(w, "q is required", http.StatusBadRequest)
		return
	}
	pages := pagesParam(r)

	var results []software.Repository
	var err error

	switch queryParam(r, "source") {
	case "gitlab":
		results, err = software.GitLabSearch(q, pages)
	case "sourceforge":
		results, err = software.SourceForgeSearch(q, pages)
	default: // github
		results, err = software.GitHubSearch(q, pages)
	}

	if err != nil {
		writeError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, results)
}

// GET /api/v1/torrents?q=...&source=piratebay|nyaa&pages=1
func handleTorrents(w http.ResponseWriter, r *http.Request) {
	q := queryParam(r, "q")
	if q == "" {
		writeError(w, "q is required", http.StatusBadRequest)
		return
	}

	switch queryParam(r, "source") {
	case "nyaa":
		results, err := torrents.NyaaSearch(q, pagesParam(r))
		if err != nil {
			writeError(w, err.Error(), http.StatusBadGateway)
			return
		}
		writeJSON(w, results)
	default: // piratebay
		results, err := torrents.PirateBaySearch(q, 0)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadGateway)
			return
		}
		writeJSON(w, results)
	}
}

// GET /api/v1/academic?q=...&source=nasa|openalex&pages=1
func handleAcademic(w http.ResponseWriter, r *http.Request) {
	q := queryParam(r, "q")
	if q == "" {
		writeError(w, "q is required", http.StatusBadRequest)
		return
	}
	pages := pagesParam(r)

	var results []academic.Paper
	var err error

	switch queryParam(r, "source") {
	case "nasa":
		results, err = academic.NASASearch(q, pages)
	default: // openalex
		results, err = academic.OpenAlexSearch(q, pages)
	}

	if err != nil {
		writeError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, results)
}

// GET /api/v1/vuln?q=...&source=nvd|exploitdb&pages=1
func handleVuln(w http.ResponseWriter, r *http.Request) {
	q := queryParam(r, "q")
	if q == "" {
		writeError(w, "q is required", http.StatusBadRequest)
		return
	}
	pages := pagesParam(r)

	switch queryParam(r, "source") {
	case "exploitdb":
		results, err := vuln.ExploitDBSearch(q, pages)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadGateway)
			return
		}
		writeJSON(w, results)
	default: // nvd
		results, err := vuln.NVDSearch(q, pages)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadGateway)
			return
		}
		writeJSON(w, results)
	}
}

// GET /api/v1/apps?q=...&source=flathub|homebrew&pages=1
func handleApps(w http.ResponseWriter, r *http.Request) {
	q := queryParam(r, "q")
	if q == "" {
		writeError(w, "q is required", http.StatusBadRequest)
		return
	}
	pages := pagesParam(r)

	var results []apps.App
	var err error

	switch queryParam(r, "source") {
	case "homebrew":
		results, err = apps.HomebrewSearch(q, pages)
	default: // flathub
		results, err = apps.FlathubSearch(q, pages)
	}

	if err != nil {
		writeError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, results)
}

// GET /api/v1/ml?q=...&source=ollama&pages=1
func handleML(w http.ResponseWriter, r *http.Request) {
	q := queryParam(r, "q")
	if q == "" {
		writeError(w, "q is required", http.StatusBadRequest)
		return
	}
	pages := pagesParam(r)

	var results []ml.Model
	var err error

	switch queryParam(r, "source") {
	case "huggingface":
		results, err = ml.HuggingFaceSearch(q, pages)
	default: // ollama
		results, err = ml.OllamaSearch(q, pages)
	}

	if err != nil {
		writeError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, results)
}

// GET /api/v1/games?q=...&source=steam&pages=1
func handleGames(w http.ResponseWriter, r *http.Request) {
	q := queryParam(r, "q")
	if q == "" {
		writeError(w, "q is required", http.StatusBadRequest)
		return
	}
	results, err := games.SteamSearch(q, pagesParam(r))
	if err != nil {
		writeError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, results)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/web", handleWeb)
	mux.HandleFunc("/api/v1/software", handleSoftware)
	mux.HandleFunc("/api/v1/torrents", handleTorrents)
	mux.HandleFunc("/api/v1/academic", handleAcademic)
	mux.HandleFunc("/api/v1/vuln", handleVuln)
	mux.HandleFunc("/api/v1/apps", handleApps)
	mux.HandleFunc("/api/v1/ml", handleML)
	mux.HandleFunc("/api/v1/games", handleGames)

	port := os.Getenv("GOSEARCH_BACKEND_PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("gosearch server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
