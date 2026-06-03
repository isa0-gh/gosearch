import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Search, Globe, Code, BookOpen, ShieldAlert, Moon, Sun, ChevronRight, LayoutGrid, Cpu, Gamepad2 } from "lucide-react";
import { doSearch } from "../api";
import type { Tab, WebResult, Repository, Paper, CVE, Exploit, App as AppType, Model as ModelType, Game, ItchGame } from "../types";
import { WebResultCard, RepoCard, PaperCard, CVECard, ExploitCard, AppCard, ModelCard, GameCard, ItchGameCard } from "./ResultCards";

const TABS: { id: Tab; icon: React.ReactNode }[] = [
  { id: "web", icon: <Globe size={14} /> },
  { id: "software", icon: <Code size={14} /> },
  { id: "apps", icon: <LayoutGrid size={14} /> },
  { id: "games", icon: <Gamepad2 size={14} /> },
  { id: "ml", icon: <Cpu size={14} /> },
  { id: "academic", icon: <BookOpen size={14} /> },
  { id: "vuln", icon: <ShieldAlert size={14} /> },
];

const SOURCE_ICONS: Record<string, string> = {
  ddg:         "https://icons.bitwarden.net/duckduckgo.com/icon.png",
  bing:        "https://icons.bitwarden.net/bing.com/icon.png",
  brave:       "https://icons.bitwarden.net/search.brave.com/icon.png",
  github:      "https://icons.bitwarden.net/github.com/icon.png",
  gitlab:      "https://icons.bitwarden.net/gitlab.com/icon.png",
  sourceforge: "https://icons.bitwarden.net/sourceforge.net/icon.png",
  openalex:    "https://icons.bitwarden.net/openalex.org/icon.png",
  nasa:        "https://icons.bitwarden.net/nasa.gov/icon.png",
  nvd:         "https://icons.bitwarden.net/nvd.nist.gov/icon.png",
  exploitdb:   "https://icons.bitwarden.net/exploit-db.com/icon.png",
  flathub:     "https://icons.bitwarden.net/flathub.org/icon.png",
  homebrew:    "https://icons.bitwarden.net/brew.sh/icon.png",
  steam:       "https://icons.bitwarden.net/store.steampowered.com/icon.png",
  itchio:      "https://icons.bitwarden.net/itch.io/icon.png",
  huggingface: "https://icons.bitwarden.net/huggingface.co/icon.png",
};

function SourcePicker({ options, value, onChange }: {
  options: string[];
  value: string;
  onChange: (v: string) => void;
}) {
  return (
    <div className="filter-group">
      <div className="source-picker">
        {options.map(opt => (
          <button
            key={opt}
            title={opt}
            className={`source-icon-btn ${value === opt ? "active" : ""}`}
            onClick={() => onChange(opt)}
          >
            <img
              src={SOURCE_ICONS[opt]}
              alt={opt}
              width={18}
              height={18}
              onError={e => { (e.currentTarget as HTMLImageElement).style.display = "none"; }}
            />
          </button>
        ))}
      </div>
    </div>
  );
}

export default function App() {
  const { t, i18n } = useTranslation();
  const [tab, setTab] = useState<Tab>("web");
  const [query, setQuery] = useState("");
  const [engine, setEngine] = useState("ddg");
  const [source, setSource] = useState("github");
  const [appsSource, setAppsSource] = useState("flathub");
  const [academicSource, setAcademicSource] = useState("openalex");
  const [gamesSource, setGamesSource] = useState("steam");
  const [vulnSource, setVulnSource] = useState("nvd");
  const [mlSource, setMlSource] = useState("ollama");
  const [page, setPage] = useState(1);
  const [results, setResults] = useState<(WebResult | Repository | Paper | CVE | Exploit | AppType | ModelType | Game | ItchGame)[]>([]);
  const [loading, setLoading] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState("");
  const [dark, setDark] = useState(() => window.matchMedia("(prefers-color-scheme: dark)").matches);

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", dark ? "dark" : "light");
  }, [dark]);

  const toggleLang = () => i18n.changeLanguage(i18n.language === "en" ? "tr" : "en");

  const fetch = async (pageNum: number, append: boolean) => {
    const src = tab === "software" ? source : tab === "academic" ? academicSource : tab === "vuln" ? vulnSource : tab === "apps" ? appsSource : tab === "ml" ? mlSource : tab === "games" ? gamesSource : source;
    const data = await doSearch(tab, query.trim(), engine, src, pageNum);
    setResults(prev => append ? [...prev, ...(data ?? [])] : (data ?? []));
    return data;
  };

  const submit = async (e?: React.FormEvent) => {
    e?.preventDefault();
    if (!query.trim()) return;
    setLoading(true);
    setError("");
    setPage(1);
    try {
      await fetch(1, false);
    } catch {
      setError(t("error"));
    } finally {
      setLoading(false);
    }
  };

  const loadMore = async () => {
    const next = page + 1;
    setLoadingMore(true);
    try {
      await fetch(next, true);
      setPage(next);
    } catch {
      setError(t("error"));
    } finally {
      setLoadingMore(false);
    }
  };

  const handleTabChange = (id: Tab) => {
    setTab(id);
    setResults([]);
    setError("");
    setPage(1);
  };

  const hasResults = results.length > 0;
  return (
    <div className="container">
      <header>
        <div className="header-row">
          <div className="logo">go<span>search</span></div>
          <div className="header-actions">
            <button className="icon-btn" onClick={() => setDark(d => !d)} aria-label="toggle theme">
              {dark ? <Sun size={15} /> : <Moon size={15} />}
            </button>
            <button className="icon-btn" onClick={toggleLang}>{t("lang")}</button>
          </div>
        </div>

        <form className="search-form" onSubmit={submit}>
          <input
            className="search-input"
            type="text"
            placeholder={t("placeholder")}
            value={query}
            onChange={e => setQuery(e.target.value)}
            autoFocus
          />
          <button className="search-btn" type="submit" disabled={loading}>
            <Search size={15} />
            <span>{t("search")}</span>
          </button>
        </form>

        <div className="tabs">
          {TABS.map(({ id, icon }) => (
            <button
              key={id}
              className={`tab-btn ${tab === id ? "active" : ""}`}
              onClick={() => handleTabChange(id)}
            >
              {icon}
              <span>{t(`tabs.${id}`)}</span>
            </button>
          ))}
        </div>

        <div className="filters">
          {tab === "web" && (
            <SourcePicker options={["ddg","bing","brave"]} value={engine} onChange={setEngine} />
          )}
          {tab === "software" && (
            <SourcePicker options={["github","gitlab","sourceforge"]} value={source} onChange={setSource} />
          )}
          {tab === "apps" && (
            <SourcePicker options={["flathub","homebrew"]} value={appsSource} onChange={setAppsSource} />
          )}
          {tab === "academic" && (
            <SourcePicker options={["openalex","nasa"]} value={academicSource} onChange={setAcademicSource} />
          )}
          {tab === "vuln" && (
            <SourcePicker options={["nvd","exploitdb"]} value={vulnSource} onChange={setVulnSource} />
          )}
          {tab === "ml" && (
            <SourcePicker options={["ollama", "huggingface"]} value={mlSource} onChange={setMlSource} />
          )}
          {tab === "games" && (
            <SourcePicker options={["steam", "itchio"]} value={gamesSource} onChange={setGamesSource} />
          )}
        </div>
      </header>

      <main>
        {loading && (
          <div className="status">
            <div className="loading-dots"><span /><span /><span /></div>
          </div>
        )}
        {error && <div className="status error">{error}</div>}
        {!loading && hasResults && (
          <div className="result-count">{results.length} {t("results")}</div>
        )}
        {!loading && !error && !hasResults && query && (
          <div className="status">{t("noResults")}</div>
        )}

        <div className="results">
          {tab === "web" && (results as WebResult[]).map((r, i) => <WebResultCard key={i} r={r} />)}
          {tab === "software" && (results as Repository[]).map((r, i) => <RepoCard key={i} r={r} />)}
          {tab === "apps" && (results as AppType[]).map((r, i) => <AppCard key={i} r={r} />)}
          {tab === "academic" && (results as Paper[]).map((r, i) => <PaperCard key={i} r={r} />)}
          {tab === "vuln" && vulnSource === "nvd" && (results as CVE[]).map((r, i) => <CVECard key={i} r={r} />)}
          {tab === "vuln" && vulnSource === "exploitdb" && (results as Exploit[]).map((r, i) => <ExploitCard key={i} r={r} />)}
          {tab === "ml" && (results as ModelType[]).map((r, i) => <ModelCard key={i} r={r} />)}
          {tab === "games" && gamesSource === "steam" && (results as Game[]).map((r, i) => <GameCard key={i} r={r} />)}
          {tab === "games" && gamesSource === "itchio" && (results as ItchGame[]).map((r, i) => <ItchGameCard key={i} r={r} />)}
        </div>

        {hasResults && (
          <footer className="load-more-footer">
            <button className="next-btn" onClick={loadMore} disabled={loadingMore}>
              {loadingMore
                ? <div className="loading-dots"><span /><span /><span /></div>
                : <><span>{t("nextPage")}</span> <ChevronRight size={15} /></>}
            </button>
          </footer>
        )}
      </main>
    </div>
  );
}
