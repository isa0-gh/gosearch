import type { Tab, WebResult, Repository, Paper, CVE, Exploit, App, Model, Game, ItchGame, GogGame } from "./types";

const BASE = "/api/v1";

export async function searchWeb(q: string, engine: string, pages: number): Promise<WebResult[]> {
  const res = await fetch(`${BASE}/web?q=${encodeURIComponent(q)}&engine=${engine}&pages=${pages}`);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function searchSoftware(q: string, source: string, pages: number): Promise<Repository[]> {
  const res = await fetch(`${BASE}/software?q=${encodeURIComponent(q)}&source=${source}&pages=${pages}`);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function searchAcademic(q: string, source: string, pages: number): Promise<Paper[]> {
  const res = await fetch(`${BASE}/academic?q=${encodeURIComponent(q)}&source=${source}&pages=${pages}`);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function searchVuln(q: string, source: string, pages: number): Promise<(CVE | Exploit)[]> {
  const res = await fetch(`${BASE}/vuln?q=${encodeURIComponent(q)}&source=${source}&pages=${pages}`);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function searchApps(q: string, source: string, pages: number): Promise<App[]> {
  const res = await fetch(`${BASE}/apps?q=${encodeURIComponent(q)}&source=${source}&pages=${pages}`);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function searchML(q: string, source: string, pages: number): Promise<Model[]> {
  const res = await fetch(`${BASE}/ml?q=${encodeURIComponent(q)}&source=${source}&pages=${pages}`);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function searchGames(q: string, source: string, pages: number): Promise<(Game | ItchGame | GogGame)[]> {
  const res = await fetch(`${BASE}/games?q=${encodeURIComponent(q)}&source=${source}&pages=${pages}`);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export function doSearch(tab: Tab, q: string, engine: string, source: string, pages: number) {
  if (tab === "web") return searchWeb(q, engine, pages);
  if (tab === "software") return searchSoftware(q, source, pages);
  if (tab === "academic") return searchAcademic(q, source, pages);
  if (tab === "vuln") return searchVuln(q, source, pages);
  if (tab === "apps") return searchApps(q, source, pages);
  if (tab === "ml") return searchML(q, source, pages);
  return searchGames(q, source, pages);
}
