export type Tab = "web" | "software" | "torrents";

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
