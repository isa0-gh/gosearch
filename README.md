# gosearch

A self-hosted metasearch engine with a Go backend and a React frontend. It aggregates results from web search engines, software repositories, academic paper databases, vulnerability databases, app stores, game stores, and LLM/ML model hubs through a unified REST API.

---

## Features

- Web search via DuckDuckGo, Bing, and Brave
- Software repository search via GitHub, GitLab, and SourceForge
- Academic paper search via OpenAlex and NASA Technical Reports Server (NTRS)
- Vulnerability search via NVD (CVE database) and Exploit-DB
- App search via Flathub and Homebrew
- Game search via Steam, itch.io, and GOG
- LLM/ML model search via Ollama and Hugging Face
- Lightweight JSON API with pagination support
- React frontend with dark mode, i18n (English and Turkish), and per-tab source switching

---

## Architecture

```
frontend/        React + Vite + TypeScript UI
cmd/server/      HTTP API server (net/http)
internal/
  scrapers/      Web search scrapers (DDG, Bing, Brave)
  software/      Repository search clients (GitHub, GitLab, SourceForge)
  academic/      Academic API clients (OpenAlex, NASA NTRS)
  vuln/          Vulnerability clients (NVD, Exploit-DB)
  apps/          Application API clients (Flathub, Homebrew)
  games/         Game store clients (Steam, itch.io, GOG)
  ml/            LLM/ML model clients (Ollama, Hugging Face)
cmd/tests/       Manual test binaries for each source
```

The frontend proxies all `/api/` requests to the backend, so only port 3000 needs to be exposed in production.

---

## API

Base URL: `http://localhost:8080/api/v1`

| Endpoint    | Parameters                                          |
|-------------|-----------------------------------------------------|
| `/web`      | `q`, `engine=ddg\|bing\|brave`, `pages`            |
| `/software` | `q`, `source=github\|gitlab\|sourceforge`, `pages` |
| `/academic` | `q`, `source=openalex\|nasa`, `pages`              |
| `/vuln`     | `q`, `source=nvd\|exploitdb`, `pages`              |
| `/apps`     | `q`, `source=flathub\|homebrew`, `pages`           |
| `/games`    | `q`, `source=steam\|itchio\|gog`, `pages`          |
| `/ml`       | `q`, `source=ollama\|huggingface`, `pages`         |

All endpoints return JSON arrays. Error responses are `{"error": "..."}` with appropriate HTTP status codes. Full schema is available in `openapi.yml`.

---

## Running with Docker

```
docker compose up --build
```

The frontend will be available at `http://localhost:3000`. The backend is not exposed directly and is only accessible through the frontend's reverse proxy.

---

## Running locally

Backend:

```
go run ./cmd/server
```

Frontend:

```
cd frontend
bun install
bun run dev
```

The Vite dev server proxies `/api` to `http://localhost:8080`.

Environment variables:

| Variable                | Default | Description       |
|-------------------------|---------|-------------------|
| `GOSEARCH_BACKEND_PORT` | `8080`  | Backend HTTP port |

---

## Testing individual sources

Each source has a standalone test binary under `cmd/tests/`:

```
go run ./cmd/tests/ddg
go run ./cmd/tests/bing
go run ./cmd/tests/brave
go run ./cmd/tests/github
go run ./cmd/tests/nasa
go run ./cmd/tests/openalex
go run ./cmd/tests/nvd
go run ./cmd/tests/exploitdb
go run ./cmd/tests/flathub
go run ./cmd/tests/homebrew
go run ./cmd/tests/steam
go run ./cmd/tests/itchio
go run ./cmd/tests/gog
go run ./cmd/tests/ollama
go run ./cmd/tests/huggingface
```

---

## Notes

Web search results rely on HTML scraping and may break if upstream sites change their markup. Academic sources (OpenAlex and NASA NTRS) use official JSON APIs and are stable. The GitHub source uses the official REST API; GitLab and SourceForge use their respective APIs as well. NVD uses the official NIST REST API v2. Exploit-DB is queried via its DataTables JSON endpoint. Flathub uses its MeiliSearch-backed API. Homebrew is queried via its official website's Algolia index. Steam and itch.io results are scraped from their public search pages. GOG is queried via its public catalog API.

---

## License

MIT
