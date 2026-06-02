import type { Tab, WebResult, Repository, Torrent, NyaaTorrent } from "./types";

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

export async function searchTorrents(q: string, source: string, pages: number): Promise<(Torrent | NyaaTorrent)[]> {
  const res = await fetch(`${BASE}/torrents?q=${encodeURIComponent(q)}&source=${source}&pages=${pages}`);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export function doSearch(tab: Tab, q: string, engine: string, source: string, pages: number) {
  if (tab === "web") return searchWeb(q, engine, pages);
  if (tab === "software") return searchSoftware(q, source, pages);
  return searchTorrents(q, source, pages);
}
