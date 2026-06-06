import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Search, Globe, Code, BookOpen, ShieldAlert, Moon, Sun, ChevronRight, LayoutGrid, Cpu, Gamepad2 } from "lucide-react";
import { doSearch } from "../api";
import type { Tab, WebResult, Repository, Paper, CVE, Exploit, App as AppType, Model as ModelType, Game, ItchGame, GogGame } from "../types";
import { WebResultCard, RepoCard, PaperCard, CVECard, ExploitCard, AppCard, ModelCard, GameCard, ItchGameCard, GogGameCard } from "./ResultCards";

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
  gog:         "https://icons.bitwarden.net/gog.com/icon.png",
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

function SearchBar({ query, setQuery, onSubmit, autoFocus }: {
  query: string;
  setQuery: (v: string) => void;
  onSubmit: (e?: React.FormEvent) => void;
  autoFocus?: boolean;
}) {
  const { t } = useTranslation();
  return (
    <form className="search-form" onSubmit={onSubmit}>
      <Search size={18} className="search-icon" />
      <input
        className="search-input"
        type="text"
        placeholder={t("placeholder")}
        value={query}
        onChange={e => setQuery(e.target.value)}
        autoFocus={autoFocus}
      />
    </form>
  );
}

const DEFAULT_SOURCES: Record<Tab, string> = {
  web: "ddg",
  software: "github",
  apps: "flathub",
  academic: "openalex",
  vuln: "nvd",
  ml: "ollama",
  games: "steam",
};

function getSourceForTab(tab: Tab, states: Record<string, string>): string {
  if (tab === "web") return states.engine;
  if (tab === "software") return states.software;
  if (tab === "apps") return states.apps;
  if (tab === "academic") return states.academic;
  if (tab === "vuln") return states.vuln;
  if (tab === "ml") return states.ml;
  return states.games;
}

function pushURL(q: string, tab: Tab, source: string) {
  const params = new URLSearchParams();
  if (q) params.set("q", q);
  if (tab !== "web") params.set("tab", tab);
  if (source && source !== DEFAULT_SOURCES[tab]) params.set("source", source);
  const qs = params.toString();
  const url = qs ? `/?${qs}` : "/";
  window.history.pushState(null, "", url);
}

function readURL(): { q: string; tab: Tab; source: string } | null {
  const params = new URLSearchParams(window.location.search);
  const q = params.get("q") ?? "";
  if (!q) return null;
  const tab = (params.get("tab") ?? "web") as Tab;
  const validTabs: Tab[] = ["web", "software", "apps", "games", "ml", "academic", "vuln"];
  if (!validTabs.includes(tab)) return { q, tab: "web", source: "ddg" };
  const source = params.get("source") ?? DEFAULT_SOURCES[tab];
  return { q, tab, source };
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
  const [results, setResults] = useState<(WebResult | Repository | Paper | CVE | Exploit | AppType | ModelType | Game | ItchGame | GogGame)[]>([]);
  const [loading, setLoading] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState("");
  const [dark, setDark] = useState(() => window.matchMedia("(prefers-color-scheme: dark)").matches);
  const [hasSearched, setHasSearched] = useState(false);

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", dark ? "dark" : "light");
  }, [dark]);

  const sourceStates = { engine, software: source, apps: appsSource, academic: academicSource, vuln: vulnSource, ml: mlSource, games: gamesSource };

  const doFetch = async (pageNum: number, append: boolean, searchTab: Tab = tab, searchQuery: string = query) => {
    const src = searchTab === "software" ? source : searchTab === "academic" ? academicSource : searchTab === "vuln" ? vulnSource : searchTab === "apps" ? appsSource : searchTab === "ml" ? mlSource : searchTab === "games" ? gamesSource : engine;
    const data = await doSearch(searchTab, searchQuery.trim(), engine, src, pageNum);
    setResults(prev => append ? [...prev, ...(data ?? [])] : (data ?? []));
    return data;
  };

  const runSearch = async (searchQuery: string, searchTab: Tab, searchSource?: string) => {
    setQuery(searchQuery);
    setTab(searchTab);
    setLoading(true);
    setError("");
    setPage(1);
    setHasSearched(true);

    if (searchSource) {
      if (searchTab === "web") setEngine(searchSource);
      else if (searchTab === "software") setSource(searchSource);
      else if (searchTab === "apps") setAppsSource(searchSource);
      else if (searchTab === "academic") setAcademicSource(searchSource);
      else if (searchTab === "vuln") setVulnSource(searchSource);
      else if (searchTab === "ml") setMlSource(searchSource);
      else if (searchTab === "games") setGamesSource(searchSource);
    }

    const src = searchSource ?? getSourceForTab(searchTab, sourceStates);
    try {
      const data = await doSearch(searchTab, searchQuery.trim(), engine, src, 1);
      setResults(data ?? []);
    } catch {
      setError(t("error"));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const initial = readURL();
    if (initial) {
      runSearch(initial.q, initial.tab, initial.source);
    }

    const onPop = () => {
      const state = readURL();
      if (state) {
        runSearch(state.q, state.tab, state.source);
      } else {
        setHasSearched(false);
        setResults([]);
        setQuery("");
        setError("");
        setPage(1);
      }
    };
    window.addEventListener("popstate", onPop);
    return () => window.removeEventListener("popstate", onPop);
  }, []);

  const toggleLang = () => i18n.changeLanguage(i18n.language === "en" ? "tr" : "en");

  const submit = async (e?: React.FormEvent) => {
    e?.preventDefault();
    if (!query.trim()) return;
    setLoading(true);
    setError("");
    setPage(1);
    setHasSearched(true);
    const src = getSourceForTab(tab, sourceStates);
    pushURL(query.trim(), tab, src);
    try {
      await doFetch(1, false);
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
      await doFetch(next, true);
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
    if (hasSearched) {
      const src = getSourceForTab(id, sourceStates);
      pushURL(query.trim(), id, src);
    }
  };

  const goHome = () => {
    setHasSearched(false);
    setResults([]);
    setQuery("");
    setError("");
    setPage(1);
    window.history.pushState(null, "", "/");
  };

  const headerActions = (
    <div className="header-actions">
      <button className="icon-btn" onClick={() => setDark(d => !d)} aria-label="toggle theme">
        {dark ? <Sun size={16} /> : <Moon size={16} />}
      </button>
      <button className="icon-btn" onClick={toggleLang}>{t("lang")}</button>
    </div>
  );

  if (!hasSearched) {
    return (
      <div className="landing">
        <div className="landing-actions">{headerActions}</div>
        <div className="landing-logo">go<span>search</span></div>
        <div className="landing-search">
          <SearchBar query={query} setQuery={setQuery} onSubmit={submit} autoFocus />
        </div>
      </div>
    );
  }

  const hasResults = results.length > 0;
  return (
    <>
      <div className="results-header">
        <div className="results-header-inner">
          <div className="results-header-top">
            <a className="header-logo" onClick={goHome}>go<span>search</span></a>
            <div className="results-search">
              <SearchBar query={query} setQuery={setQuery} onSubmit={submit} />
            </div>
            {headerActions}
          </div>

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
        </div>
      </div>

      <div className="filters">
        {tab === "web" && <SourcePicker options={["ddg","bing","brave"]} value={engine} onChange={setEngine} />}
        {tab === "software" && <SourcePicker options={["github","gitlab","sourceforge"]} value={source} onChange={setSource} />}
        {tab === "apps" && <SourcePicker options={["flathub","homebrew"]} value={appsSource} onChange={setAppsSource} />}
        {tab === "academic" && <SourcePicker options={["openalex","nasa"]} value={academicSource} onChange={setAcademicSource} />}
        {tab === "vuln" && <SourcePicker options={["nvd","exploitdb"]} value={vulnSource} onChange={setVulnSource} />}
        {tab === "ml" && <SourcePicker options={["ollama", "huggingface"]} value={mlSource} onChange={setMlSource} />}
        {tab === "games" && <SourcePicker options={["steam", "itchio", "gog"]} value={gamesSource} onChange={setGamesSource} />}
      </div>

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
          {tab === "games" && gamesSource === "gog" && (results as GogGame[]).map((r, i) => <GogGameCard key={r.ID || i} r={r} />)}
        </div>

        {hasResults && (
          <div className="load-more-footer">
            <button className="next-btn" onClick={loadMore} disabled={loadingMore}>
              {loadingMore
                ? <div className="loading-dots"><span /><span /><span /></div>
                : <><span>{t("nextPage")}</span> <ChevronRight size={15} /></>}
            </button>
          </div>
        )}
      </main>
    </>
  );
}
