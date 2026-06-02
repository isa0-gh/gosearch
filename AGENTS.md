# AI Agent Onboarding & Architecture Manual (AGENTS.md)

Welcome, AI Developer! This document is your single source of truth for understanding, modifying, and extending the **gosearch** codebase. It outlines the codebase architecture, design patterns, engineering standards, and provides a step-by-step cookbook for implementing new search sources.

---

## Architecture Blueprint

**gosearch** is a self-hosted metasearch engine consisting of a lightweight Go backend API and a modern React + TypeScript frontend.

```
                  ┌──────────────────────────────────────┐
                  │          React Frontend UI           │
                  │    (Vite + TS + i18next + CSS)       │
                  └──────────────────┬───────────────────┘
                                     │ (Fetch API Proxy /api/v1/*)
                                     ▼
                  ┌──────────────────────────────────────┐
                  │              Go Backend              │
                  │  (net/http standard multiplexer)     │
                  └──────────────────┬───────────────────┘
                                     │
           ┌─────────────────────────┼─────────────────────────┐
           ▼                         ▼                         ▼
 ┌───────────────────┐     ┌───────────────────┐     ┌───────────────────┐
 │   Web Scrapers    │     │   API Clients     │     │ Academic/Vuln/etc │
 │ (DDG, Bing, Brave)│     │(GitHub, Flathub)  │     │(OpenAlex, NVD...) │
 └───────────────────┘     └───────────────────┘     └───────────────────┘
```

### Directory Structure Map
- `/cmd/server/`: The backend HTTP server entry point (`main.go`). It binds handlers and manages query dispatching.
- `/cmd/tests/`: Independent test binaries for each scraping engine/API client. Extremely valuable for rapid debugging.
- `/internal/`: Categorized packages containing client/scraper business logic:
  - `academic/`: Clients for research databases (OpenAlex, NASA NTRS).
  - `apps/`: Search engines for application packages (Flathub, Homebrew).
  - `scrapers/`: Raw HTML scrapers for standard web search (DuckDuckGo, Bing, Brave).
  - `software/`: Software repository indexers (GitHub, GitLab, SourceForge).
  - `torrents/`: Index scrapers (The Pirate Bay, Nyaa).
  - `vuln/`: Security vulnerability indices (NVD, Exploit-DB).
- `/frontend/`: The React client interface:
  - `src/components/`: Modular UI, mainly `App.tsx` (state & layout) and `ResultCards.tsx` (card rendering).
  - `src/api.ts`: API clients querying Go HTTP endpoints.
  - `src/types.ts`: TypeScript data models matching Go structs.
  - `src/i18n.ts`: English and Turkish translation keys.

---

## Engineering Standards & Coding Style

To maintain the clean, low-dependency design of this codebase, you must adhere to the following rules:

### 1. Go Backend Guidelines
- **Zero External Routers Policy**: Avoid importing frameworks like Gin, Echo, or Chi. Use Go's standard library `net/http` and `http.NewServeMux` exclusively.
- **Strict Parsing Precautions**:
  - Always validate search queries: if `q == ""`, return a `http.StatusBadRequest` with a JSON payload `{"error": "q is required"}`.
  - Always default page parameters cleanly. Use the `pagesParam` helper to gracefully return `1` on missing or malformed inputs.
- **Resource Cleanup**: Always close response bodies:
  ```go
  resp, err := client.Do(req)
  if err != nil {
      return nil, err
  }
  defer resp.Body.Close()
  ```
- **Error Wrapping**: Wrap HTTP client errors with clear context (e.g. `fmt.Errorf("page %d: %w", p, err)`), and bubble errors up to the HTTP handler to respond with `http.StatusBadGateway`.
- **Minimal Struct Footprints**: For external JSON API requests, declare anonymous/in-line structs inside your search function to avoid namespace pollution in your package unless the schema is shared.

### 2. Web Scraping Guidelines
- **Resilient goquery Selectors**: When scraping HTML (e.g. using `github.com/PuerkitoBio/goquery`), protect against upstream DOM changes by gracefully skipping empty nodes instead of crashing.
- **Strict Headers**: Upstream sites will reject requests without realistic user headers. Always use standard browser headers:
  ```go
  req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")
  ```
- **Url Query Escaping**: Always encode client query input using `url.QueryEscape(query)` or `url.Values{}` to avoid injection and broken request links.

### 3. Frontend React & TypeScript Guidelines
- **Translation First**: Do not hardcode English or Turkish user interface labels. Always register fields inside `/frontend/src/i18n.ts` and call the translator hook:
  ```tsx
  const { t } = useTranslation();
  // usage: {t("fields.updated")}
  ```
- **Strict TypeScript Typing**: Avoid the use of `any` types. When introducing new results, declare complete schemas inside `types.ts` and use them as generic types or specific parameters in `ResultCards.tsx`.
- **CSS-Variable Centric Styling**: Support light and dark modes natively. Use defined CSS variables (like `var(--border)`, `var(--bg-secondary)`) in component styles instead of arbitrary hex colors.

### 4. Git Commit Standards (Conventional Commits)
This repository strictly adheres to the Conventional Commits specification. Commit messages must use lowercase categories and specify exact scopes where applicable. 

Format: `type(scope): description` or `type: description`

Common categories used in this codebase:
- **feat**: New user-facing features or system capabilities (e.g., `feat(frontend): add apps tab with source picker`, `feat(apps): add homebrew search module`).
- **fix**: Bug fixes (e.g., `fix(apps): use 1-based indexing for flathub pagination`).
- **docs**: Documentation modifications (e.g., `docs(openapi): add apps endpoint`, `docs: update README with architecture details`).
- **test**: Standalone test programs, verification scripts, or unit tests (e.g., `test(apps): add flathub search test script`).

---

## The "Adding a New Search Provider" Step-by-Step Playbook

Follow this precise sequence to integrate a new search source:

### Step 1: Create the Backend Provider client
Create a Go client inside `/internal/<category>/<source>.go`.
Here is a standardized structural template:

```go
package category // e.g. software, academic, vuln, etc.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type NewItem struct {
	Title       string
	URL         string
	Description string
}

func NewSourceSearch(query string, pages int) ([]NewItem, error) {
	client := &http.Client{}
	var results []NewItem

	for p := 1; p <= pages; p++ {
		escaped := url.QueryEscape(query)
		u := fmt.Sprintf("https://api.example.com/search?q=%s&page=%d", escaped, p)
		
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return results, err
		}
		req.Header.Set("User-Agent", "gosearch")

		resp, err := client.Do(req)
		if err != nil {
			return results, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return results, fmt.Errorf("example API status: %d", resp.StatusCode)
		}

		var apiResponse struct {
			Data []struct {
				Name string `json:"name"`
				Link string `json:"link"`
				Desc string `json:"description"`
			} `json:"data"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			resp.Body.Close()
			return results, err
		}
		resp.Body.Close()

		if len(apiResponse.Data) == 0 {
			break
		}

		for _, item := range apiResponse.Data {
			results = append(results, NewItem{
				Title:       item.Name,
				URL:         item.Link,
				Description: item.Desc,
			})
		}
	}

	return results, nil
}
```

### Step 2: Register API Server Routing
Open `/cmd/server/main.go`:
1. Import the package if adding a new directory category.
2. Update the corresponding `handle<Category>` method (e.g. `handleSoftware`, `handleAcademic`) to introduce the new source case.
3. If this source is a completely new category, create a new handler `handleNewCategory` and register it with `mux.HandleFunc("/api/v1/newcategory", handleNewCategory)` in `main()`.

### Step 3: Document API in OpenAPI Contract
Open `/openapi.yml`:
1. Find the path (e.g., `/software`, `/academic`) and update the `source` parameter's `enum` list with your new lowercase key.
2. If you added a new schema structure, declare it under `components/schemas/` to keep our API fully documented.

### Step 4: Add Independent Test Binary
Create a standalone testing main function under `/cmd/tests/<source>/main.go`:
```go
package main

import (
	"fmt"
	"log"

	"github.com/isa0-gh/gosearch/internal/category"
)

func main() {
	results, err := category.NewSourceSearch("test query", 1)
	if err != nil {
		log.Fatal(err)
	}
	for i, r := range results {
		fmt.Printf("[%d] %s\n    %s\n\n", i+1, r.Title, r.URL)
	}
}
```
> [!TIP]
> Execute `go run ./cmd/tests/<source>` immediately to isolate scraper errors before launching the web server.

### Step 5: Declare Frontend Model Types
Open `/frontend/src/types.ts`:
1. Define a TypeScript interface that matches the JSON keys returned by the backend.
2. Ensure capitalization matches the exported fields of your Go struct:
   ```typescript
   export interface NewItem {
     Title: string;
     URL: string;
     Description: string;
   }
   ```

### Step 6: Create Frontend Fetch Call
Open `/frontend/src/api.ts`:
1. Import the new type from `./types`.
2. Ensure the query handler maps to your REST endpoint.

### Step 7: Build UI Card Component
Open `/frontend/src/components/ResultCards.tsx`:
1. Create a specialized rendering card:
   ```tsx
   export function NewItemCard({ r }: { r: NewItem }) {
     return (
       <div className="result-item">
         <div className="result-title">
           <a href={r.URL} target="_blank" rel="noopener noreferrer">{r.Title}</a>
         </div>
         {r.Description && <div className="result-snippet">{r.Description}</div>}
       </div>
     );
   }
   ```
2. Export your card and ensure it's imported correctly.

### Step 8: Update Translations
Open `/frontend/src/i18n.ts`:
1. Add localized strings for any fields or categories you introduced.

### Step 9: Wire Component into Layout
Open `/frontend/src/components/App.tsx`:
1. Import your brand new Card.
2. Add your engine/source string options inside the engine selectors / state handlers.
3. Within the result render loop, map the search category to the new layout card:
   ```tsx
   {tab === "newcategory" && results.map((r, i) => (
     <NewItemCard key={i} r={r as NewItem} />
   ))}
   ```

---

## Sandbox Execution Commands

Here are the critical terminal commands you will need:

### Start Backend Locally
```bash
go run ./cmd/server
```
*Listens by default on `http://localhost:8080`*

### Start Frontend Dev Server
```bash
cd frontend
bun run dev
```
*Listens on `http://localhost:3000` and proxies `/api` to the backend.*

### Verify API / Run Test Binary
```bash
go run ./cmd/tests/<source>
```

### Complete Stack Docker Build
```bash
docker compose up --build
```

---

## Common Trapdoors & Troubleshooting Checklist

- **HTML Class Alterations**: Public search engines frequently change DOM markup to avoid scraping. If a web search scraper breaks:
  - Run the test binary locally and write the HTML output to a scratch file to inspect selector accuracy.
  - Utilize `goquery` fallback checks.
- **Rate-Limits & IP Blocks**: If your request returns HTTP status `429 Too Many Requests` or `403 Forbidden`:
  - Verify headers. Some endpoints (e.g. Exploit-DB, NVD) reject standard Go clients immediately unless header formats strictly match a normal browser.
- **Pagination Inconsistencies**: Some search sources use `0-indexed` pagination (e.g. page parameter starting at 0) while others are `1-indexed`. Check API documentation and adjust internal Go loops accordingly to prevent off-by-one missing result pages.
- **Date String Anomalies**: Go standard library json decoder handles `time.Time` automatically *only* if the source format complies strictly with RFC3339. Otherwise, decode the date fields as raw string or int64, and parse/convert them cleanly before returning.

Now you are fully armed to develop with **gosearch**. Good luck!
