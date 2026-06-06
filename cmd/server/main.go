package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/isa0-gh/gosearch/internal/academic"
	"github.com/isa0-gh/gosearch/internal/apps"
	"github.com/isa0-gh/gosearch/internal/games"
	"github.com/isa0-gh/gosearch/internal/ml"
	"github.com/isa0-gh/gosearch/internal/scrapers"
	"github.com/isa0-gh/gosearch/internal/software"
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

// GET /api/v1/games?q=...&source=steam|itchio|gog&pages=1
func handleGames(w http.ResponseWriter, r *http.Request) {
	q := queryParam(r, "q")
	if q == "" {
		writeError(w, "q is required", http.StatusBadRequest)
		return
	}
	pages := pagesParam(r)

	switch queryParam(r, "source") {
	case "itchio":
		results, err := games.ItchSearch(q, pages)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadGateway)
			return
		}
		writeJSON(w, results)
	case "gog":
		results, err := games.GogSearch(q, pages)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadGateway)
			return
		}
		writeJSON(w, results)
	default: // steam
		results, err := games.SteamSearch(q, pages)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadGateway)
			return
		}
		writeJSON(w, results)
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

type logEntry struct {
	Timestamp  string `json:"timestamp"`
	IP         string `json:"ip"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Query      string `json:"query,omitempty"`
	Status     int    `json:"status"`
	Duration   string `json:"duration"`
	UserAgent  string `json:"user_agent,omitempty"`
	Referer    string `json:"referer,omitempty"`
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return splitFirst(xff, ',')
	}
	if xri := r.Header.Get("X-Real-Ip"); xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func splitFirst(s string, sep byte) string {
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			return strings.TrimSpace(s[:i])
		}
	}
	return strings.TrimSpace(s)
}

func loggerMiddleware(next http.Handler) http.Handler {
	enc := json.NewEncoder(os.Stdout)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)

		entry := logEntry{
			Timestamp: start.UTC().Format(time.RFC3339),
			IP:        clientIP(r),
			Method:    r.Method,
			Path:      r.URL.Path,
			Query:     r.URL.RawQuery,
			Status:    sw.status,
			Duration:  time.Since(start).String(),
			UserAgent: r.UserAgent(),
			Referer:   r.Referer(),
		}
		enc.Encode(entry)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/web", handleWeb)
	mux.HandleFunc("/api/v1/software", handleSoftware)
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
	log.Fatal(http.ListenAndServe(addr, loggerMiddleware(mux)))
}
