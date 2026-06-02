# gosearch

A self-hosted metasearch engine with a Go backend and a React frontend. It aggregates results from web search engines, software repositories, torrent indexes, academic paper databases, and vulnerability databases through a unified REST API.

---

## Features

- Web search via DuckDuckGo, Bing, and Brave
- Software repository search via GitHub, GitLab, and SourceForge
- Torrent search via The Pirate Bay and Nyaa
- Academic paper search via OpenAlex and NASA Technical Reports Server (NTRS)
- Vulnerability search via NVD (CVE database) and Exploit-DB
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
  torrents/      Torrent scrapers (Pirate Bay, Nyaa)
  academic/      Academic API clients (OpenAlex, NASA NTRS)
  vuln/          Vulnerability clients (NVD, Exploit-DB)
cmd/tests/       Manual test binaries for each source
```

The frontend proxies all `/api/` requests to the backend, so only port 3000 needs to be exposed in production.

---

## API

Base URL: `http://localhost:8080/api/v1`

| Endpoint      | Parameters                                              |
|---------------|---------------------------------------------------------|
| `/web`        | `q`, `engine=ddg\|bing\|brave`, `pages`                |
| `/software`   | `q`, `source=github\|gitlab\|sourceforge`, `pages`     |
| `/torrents`   | `q`, `source=piratebay\|nyaa`, `pages`                 |
| `/academic`   | `q`, `source=openalex\|nasa`, `pages`                  |
| `/vuln`       | `q`, `source=nvd\|exploitdb`, `pages`                  |

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
go run ./cmd/tests/piratebay
go run ./cmd/tests/nyaa
go run ./cmd/tests/nasa
go run ./cmd/tests/openalex
go run ./cmd/tests/nvd
go run ./cmd/tests/exploitdb
```

---

## Notes

Web search results rely on HTML scraping and may break if upstream sites change their markup. Academic sources (OpenAlex and NASA NTRS) use official JSON APIs and are stable. The GitHub source uses the official REST API; GitLab and SourceForge use their respective APIs as well. NVD uses the official NIST REST API v2. Exploit-DB is queried via its DataTables JSON endpoint.

---

## License

MIT
