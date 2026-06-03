import { useTranslation } from "react-i18next";
import { Star, Globe, Download, Package, Copy, Check, Monitor, Apple, Terminal } from "lucide-react";
import { useState, type ReactNode } from "react";
import type { WebResult, Repository, Paper, CVE, Exploit, App, Model, Game, ItchGame } from "../types";

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
  const [copied, setCopied] = useState(false);

  const isHomebrew = r.URL.includes("formulae.brew.sh");
  const isCask = r.Developer === "Homebrew Cask";
  const installCmd = isHomebrew
    ? `brew install ${isCask ? "--cask " : ""}${r.AppID}`
    : `flatpak install flathub ${r.AppID}`;

  const copy = () => {
    navigator.clipboard.writeText(installCmd);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="result-item">
      <div className="result-title" style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          {r.Icon ? (
             <img src={r.Icon} alt="" width={24} height={24} onError={e => e.currentTarget.style.display='none'} />
          ) : <Package size={20} />}
          <a href={r.URL} target="_blank" rel="noopener noreferrer">{r.Name}</a>
        </div>
        <button 
          className="copy-command-btn" 
          onClick={copy} 
          title={installCmd}
          style={{ 
            display: 'flex', 
            alignItems: 'center', 
            gap: '4px', 
            fontSize: '11px', 
            padding: '4px 8px',
            borderRadius: '4px',
            border: '1px solid var(--border)',
            background: 'var(--bg-secondary)',
            color: 'var(--text-secondary)',
            cursor: 'pointer'
          }}
        >
          {copied ? <Check size={12} color="#4caf50" /> : <Copy size={12} />}
          <code>{installCmd}</code>
        </button>
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

export function ModelCard({ r }: { r: Model }) {  const { t } = useTranslation();
  return (
    <div className="result-item">
      <div className="result-title">
        <a href={r.url} target="_blank" rel="noopener noreferrer">{r.name}</a>
      </div>
      {r.description && <div className="result-desc">{r.description}</div>}
      {r.capabilities && r.capabilities.length > 0 && (
        <div style={{ display: "flex", gap: "6px", flexWrap: "wrap", marginTop: "6px", marginBottom: "6px" }}>
          {r.capabilities.map((cap) => {
            let bg = "var(--bg-secondary)";
            let fg = "var(--text-secondary)";
            if (cap === "vision") {
              bg = "rgba(139, 92, 246, 0.1)";
              fg = "rgb(139, 92, 246)";
            } else if (cap === "tools") {
              bg = "rgba(16, 185, 129, 0.1)";
              fg = "rgb(16, 185, 129)";
            } else if (cap === "thinking") {
              bg = "rgba(59, 130, 246, 0.1)";
              fg = "rgb(59, 130, 246)";
            } else if (cap === "cloud") {
              bg = "rgba(6, 182, 212, 0.1)";
              fg = "rgb(6, 182, 212)";
            }
            return (
              <span
                key={cap}
                className="meta-pill"
                style={{
                  backgroundColor: bg,
                  color: fg,
                  borderColor: fg,
                  fontSize: "11px",
                  padding: "2px 6px",
                  borderRadius: "4px",
                  fontWeight: 500,
                }}
              >
                {cap}
              </span>
            );
          })}
        </div>
      )}
      <div className="result-meta">
        {r.pulls && (
          <span style={{ display: "flex", alignItems: "center", gap: "4px" }}>
            <Download size={11} /> {r.pulls} {t("fields.pulls", "Pulls")}
          </span>
        )}
        {r.tags && (
          <span style={{ display: "flex", alignItems: "center", gap: "4px" }}>
            <Package size={11} /> {r.tags} {t("fields.tags", "Tags")}
          </span>
        )}
        {r.size && (
          <span className="meta-pill" style={{ color: "var(--text-primary)", borderColor: "var(--border)" }}>
            {r.size}
          </span>
        )}
        {r.updated && (
          <span>
            {t("fields.updated")}: {r.updated}
          </span>
        )}
      </div>
    </div>
  );
}




const PLATFORM_ICONS: Record<string, ReactNode> = {
  win:   <Monitor size={13} />,
  mac:   <Apple size={13} />,
  linux: <Terminal size={13} />,
};

const REVIEW_COLOR: Record<string, string> = {
  positive:          "#4caf50",
  mixed:             "#f57c00",
  negative:          "#e53935",
  overwhelminglyPos: "#2e7d32",
  overwhelminglyNeg: "#b71c1c",
  mostlyPositive:    "#66bb6a",
  mostlyNegative:    "#ef5350",
};

export function GameCard({ r }: { r: Game }) {
  const reviewColor = r.ReviewClass ? (REVIEW_COLOR[r.ReviewClass] ?? "var(--muted)") : undefined;
  return (
    <div className="result-item" style={{ display: "flex", gap: "12px" }}>
      {r.ImageURL && (
        <img
          src={r.ImageURL}
          alt=""
          width={92}
          height={43}
          style={{ borderRadius: "4px", objectFit: "cover", flexShrink: 0 }}
          onError={e => (e.currentTarget.style.display = "none")}
        />
      )}
      <div style={{ flex: 1, minWidth: 0 }}>
        <div className="result-title">
          <a href={r.URL} target="_blank" rel="noopener noreferrer">{r.Title}</a>
        </div>
        {r.ReviewSummary && (
          <div className="result-snippet" style={{ color: reviewColor, fontSize: "12px" }}>
            {r.ReviewSummary.replace(/&lt;br&gt;/g, " ").replace(/&amp;/g, "&")}
          </div>
        )}
        <div className="result-meta">
          {r.Price && (
            <span style={{ fontWeight: 600 }}>
              {r.DiscountPercent && (
                <span style={{ color: "#4caf50", marginRight: 4 }}>{r.DiscountPercent}</span>
              )}
              {r.OriginalPrice && (
                <span style={{ textDecoration: "line-through", color: "var(--muted)", marginRight: 4, fontWeight: 400 }}>
                  {r.OriginalPrice}
                </span>
              )}
              {r.Price}
            </span>
          )}
          {r.ReleaseDate && <span>{r.ReleaseDate}</span>}
          {r.Platforms?.map(p => (
            <span key={p} title={p}>{PLATFORM_ICONS[p] ?? p}</span>
          ))}
        </div>
      </div>
    </div>
  );
}

export function ItchGameCard({ r }: { r: ItchGame }) {
  return (
    <div className="result-item" style={{ display: "flex", gap: "12px" }}>
      {r.ThumbnailURL && (
        <img
          src={r.ThumbnailURL}
          alt=""
          width={92}
          height={69}
          style={{ borderRadius: "4px", objectFit: "cover", flexShrink: 0 }}
          onError={e => (e.currentTarget.style.display = "none")}
        />
      )}
      <div style={{ flex: 1, minWidth: 0 }}>
        <div className="result-title">
          <a href={r.URL} target="_blank" rel="noopener noreferrer">{r.Title}</a>
        </div>
        {r.Author && (
          <div className="result-url">
            <a href={r.AuthorURL} target="_blank" rel="noopener noreferrer">{r.Author}</a>
          </div>
        )}
        {r.Description && <div className="result-snippet">{r.Description}</div>}
        <div className="result-meta">
          {r.Rating && (
            <span title={`${r.Rating.Total} ratings`}>
              ★ {r.Rating.Average.toFixed(2)}
              <span style={{ color: "var(--muted)", marginLeft: 3 }}>({r.Rating.Total})</span>
            </span>
          )}
          {r.Genre && <span className="meta-pill">{r.Genre}</span>}
          {r.Platforms.Windows && <span title="Windows"><Monitor size={13} /></span>}
          {r.Platforms.MacOS && <span title="macOS"><Apple size={13} /></span>}
          {r.Platforms.Linux && <span title="Linux"><Terminal size={13} /></span>}
          {r.Platforms.Web && <span title="Web"><Globe size={13} /></span>}
          {r.Platforms.Android && <span title="Android"><Download size={13} /></span>}
        </div>
      </div>
    </div>
  );
}
