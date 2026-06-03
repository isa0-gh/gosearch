export type Tab = "web" | "software" | "torrents" | "academic" | "vuln" | "apps" | "ml" | "games";

export interface Game {
  AppID: string;
  Title: string;
  URL: string;
  ImageURL: string;
  ReleaseDate: string;
  Price: string;
  OriginalPrice: string;
  DiscountPercent: string;
  ReviewSummary: string;
  ReviewClass: string;
  Platforms: string[];
}

export interface Model {
  name: string;
  url: string;
  description: string;
  capabilities: string[];
  pulls: string;
  tags: string;
  size: string;
  updated: string;
}

export interface App {
  AppID: string;
  Name: string;
  Summary: string;
  Developer: string;
  License: string;
  Icon: string;
  URL: string;
  UpdatedAt: number;
}

export interface CVE {
  ID: string;
  Description: string;
  Published: string;
  Severity: string;
  Score: number;
  URL: string;
}

export interface Exploit {
  ID: string;
  Title: string;
  Type: string;
  Platform: string;
  Author: string;
  Published: string;
  CVEs: string[];
  URL: string;
}

export interface Paper {
  Title: string;
  URL: string;
  Authors: string;
  Abstract: string;
  Type: string;
}

export interface WebResult {
  Title: string;
  URL: string;
  Snippet: string;
}

export interface Repository {
  Name: string;
  URL: string;
  Description: string;
  Stars: number;
  Language: string;
  UpdatedAt: string;
}

export interface Torrent {
  Name: string;
  InfoHash: string;
  MagnetURL: string;
  Seeders: number;
  Leechers: number;
  Size: number;
  AddedAt: string;
  Category: string;
  Uploader: string;
}

export interface NyaaTorrent {
  Name: string;
  URL: string;
  MagnetURL: string;
  Size: string;
  Category: string;
  Seeders: number;
  Leechers: number;
  Downloads: number;
  AddedAt: string;
}
