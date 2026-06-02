import { useTranslation } from "react-i18next";
import { Star, Globe, Download, Magnet, Package } from "lucide-react";
import type { WebResult, Repository, Torrent, NyaaTorrent, Paper, CVE, Exploit, App } from "../types";

function favicon(url: string) {
  try {
    return `https://icons.bitwarden.net/${new URL(url).hostname}/icon.png`;
  } catch {
    return null;
  }
}

export function WebResultCard({ r }: { r: WebResult }) {
  const icon = favicon(r.URL);
  return (
    <div className="result-item">
      <div className="result-url">
        {icon && (
          <img
            src={icon}
            alt=""
            width={13}
            height={13}
            style={{ marginRight: 5, verticalAlign: "middle", opacity: 0.8 }}
            onError={e => (e.currentTarget.style.display = "none")}
          />
        )}
        {r.URL}
      </div>
      <div className="result-title">
        <a href={r.URL} target="_blank" rel="noopener noreferrer">{r.Title}</a>
      </div>
      <div className="result-snippet">{r.Snippet}</div>
    </div>
  );
}

export function RepoCard({ r }: { r: Repository }) {
  const { t } = useTranslation();
  return (
    <div className="result-item">
      <div className="result-title">
        <a href={r.URL} target="_blank" rel="noopener noreferrer">{r.Name}</a>
      </div>
      {r.Description && <div className="result-desc">{r.Description}</div>}
      <div className="result-meta">
        {r.Stars > 0 && <span><Star size={11} /> {r.Stars.toLocaleString()}</span>}
        {r.Language && <span className="meta-pill">{r.Language}</span>}
        {r.UpdatedAt && <span>{t("fields.updated")} {new Date(r.UpdatedAt).toLocaleDateString()}</span>}
      </div>
    </div>
  );
}

export function TorrentCard({ r }: { r: Torrent | NyaaTorrent }) {
  const { t } = useTranslation();
  const isPirate = "InfoHash" in r;
  return (
    <div className="result-item">
      <div className="result-title">
        {"URL" in r
          ? <a href={(r as NyaaTorrent).URL} target="_blank" rel="noopener noreferrer">{r.Name}</a>
          : r.Name}
      </div>
      <div className="result-meta">
        <span style={{ color: "#4caf50" }}>↑ {r.Seeders}</span>
        <span style={{ color: "#e57373" }}>↓ {r.Leechers}</span>
        {"Downloads" in r && <span><Download size={11} /> {(r as NyaaTorrent).Downloads}</span>}
        {isPirate && (r as Torrent).Size > 0 && (
          <span>{((r as Torrent).Size / 1_073_741_824) >= 1
            ? `${((r as Torrent).Size / 1_073_741_824).toFixed(2)} GB`
            : `${((r as Torrent).Size / 1_048_576).toFixed(1)} MB`}
          </span>
        )}
        {!isPirate && (r as NyaaTorrent).Size && <span>{(r as NyaaTorrent).Size}</span>}
        {r.Category && <span className="meta-pill">{r.Category}</span>}
        {isPirate && (r as Torrent).Uploader && <span>{t("fields.uploader")}: {(r as Torrent).Uploader}</span>}
      </div>
      {r.MagnetURL && (
        <a className="magnet-link" href={r.MagnetURL}>
          <Magnet size={11} /> magnet
        </a>
      )}
    </div>
  );
}

export function PaperCard({ r }: { r: Paper }) {
  return (
    <div className="result-item">
      <div className="result-title">
        {r.URL
          ? <a href={r.URL} target="_blank" rel="noopener noreferrer">{r.Title}</a>
          : r.Title}
      </div>
      {r.Authors && <div className="result-url">{r.Authors}</div>}
      {r.Abstract && <div className="result-snippet">{r.Abstract}</div>}
      {r.Type && <div className="result-meta"><span className="meta-pill">{r.Type}</span></div>}
    </div>
  );
}

const SEVERITY_COLOR: Record<string, string> = {
  CRITICAL: "#d32f2f",
  HIGH:     "#e64a19",
  MEDIUM:   "#f57c00",
  LOW:      "#388e3c",
};

export function CVECard({ r }: { r: CVE }) {
  const color = SEVERITY_COLOR[r.Severity?.toUpperCase()] ?? "var(--muted)";
  return (
    <div className="result-item">
      <div className="result-title">
        <a href={r.URL} target="_blank" rel="noopener noreferrer">{r.ID}</a>
      </div>
      {r.Description && <div className="result-snippet">{r.Description}</div>}
      <div className="result-meta">
        {r.Severity && <span className="meta-pill" style={{ color, borderColor: color }}>{r.Severity}</span>}
        {r.Score > 0 && <span>CVSS {r.Score.toFixed(1)}</span>}
        {r.Published && <span>{r.Published.slice(0, 10)}</span>}
      </div>
    </div>
  );
}

export function ExploitCard({ r }: { r: Exploit }) {
  return (
    <div className="result-item">
      <div className="result-title">
        <a href={r.URL} target="_blank" rel="noopener noreferrer">{r.Title}</a>
      </div>
      <div className="result-meta">
        {r.Type && <span className="meta-pill">{r.Type}</span>}
        {r.Platform && <span className="meta-pill">{r.Platform}</span>}
        {r.Author && <span>{r.Author}</span>}
        {r.Published && <span>{r.Published}</span>}
        {r.CVEs?.map(c => <span key={c} className="meta-pill">{c}</span>)}
      </div>
    </div>
  );
}

export function AppCard({ r }: { r: App }) {
  const { t } = useTranslation();
  return (
    <div className="result-item">
      <div className="result-title" style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
        {r.Icon ? (
           <img src={r.Icon} alt="" width={24} height={24} onError={e => e.currentTarget.style.display='none'} />
        ) : <Package size={20} />}
        <a href={r.URL} target="_blank" rel="noopener noreferrer">{r.Name}</a>
      </div>
      {r.Summary && <div className="result-desc">{r.Summary}</div>}
      <div className="result-meta">
        {r.Developer && <span>{r.Developer}</span>}
        {r.License && <span className="meta-pill">{r.License}</span>}
        {r.UpdatedAt > 0 && <span>{t("fields.updated")} {new Date(r.UpdatedAt * 1000).toLocaleDateString()}</span>}
      </div>
    </div>
  );
}

