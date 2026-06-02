import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Search, Globe, Code, Download, Moon, Sun, ChevronRight } from "lucide-react";
import { doSearch } from "../api";
import type { Tab, WebResult, Repository, Torrent, NyaaTorrent } from "../types";
import { WebResultCard, RepoCard, TorrentCard } from "./ResultCards";

const TABS: { id: Tab; icon: React.ReactNode }[] = [
  { id: "web", icon: <Globe size={14} /> },
  { id: "software", icon: <Code size={14} /> },
  { id: "torrents", icon: <Download size={14} /> },
];

export default function App() {
  const { t, i18n } = useTranslation();
  const [tab, setTab] = useState<Tab>("web");
  const [query, setQuery] = useState("");
  const [engine, setEngine] = useState("ddg");
  const [source, setSource] = useState("github");
  const [torrentSource, setTorrentSource] = useState("piratebay");
  const [page, setPage] = useState(1);
  const [results, setResults] = useState<(WebResult | Repository | Torrent | NyaaTorrent)[]>([]);
  const [loading, setLoading] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState("");
  const [dark, setDark] = useState(() => window.matchMedia("(prefers-color-scheme: dark)").matches);

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", dark ? "dark" : "light");
  }, [dark]);

  const toggleLang = () => i18n.changeLanguage(i18n.language === "en" ? "tr" : "en");

  const fetch = async (pageNum: number, append: boolean) => {
    const src = tab === "software" ? source : torrentSource;
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
  // piratebay returns all results in one shot — no next page
  const canLoadMore = tab !== "torrents" || torrentSource === "nyaa";

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
            <div className="filter-group">
              <label>{t("engine")}</label>
              <select value={engine} onChange={e => setEngine(e.target.value)}>
                <option value="ddg">DuckDuckGo</option>
                <option value="bing">Bing</option>
                <option value="brave">Brave</option>
              </select>
            </div>
          )}
          {tab === "software" && (
            <div className="filter-group">
              <label>{t("source")}</label>
              <select value={source} onChange={e => setSource(e.target.value)}>
                <option value="github">GitHub</option>
                <option value="gitlab">GitLab</option>
                <option value="sourceforge">SourceForge</option>
              </select>
            </div>
          )}
          {tab === "torrents" && (
            <div className="filter-group">
              <label>{t("source")}</label>
              <select value={torrentSource} onChange={e => setTorrentSource(e.target.value)}>
                <option value="piratebay">Pirate Bay</option>
                <option value="nyaa">Nyaa</option>
              </select>
            </div>
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
          {tab === "torrents" && (results as (Torrent | NyaaTorrent)[]).map((r, i) => <TorrentCard key={i} r={r} />)}
        </div>

        {hasResults && canLoadMore && (
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
