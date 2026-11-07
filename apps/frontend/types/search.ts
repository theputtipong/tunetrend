export type SearchSource = "chart" | "searched" | "youtube";

export interface SearchItem {
  id: string;
  title: string;
  channelTitle: string;
  thumbnailUrl: string;
  viewCount: string;
  publishedAt: string;
  videoType: string;
  categoryId: string;
  source: SearchSource;
  chartTab?: "trending" | "new";
  chartRank?: number;
}

export interface SearchedQuery {
  query: string;
  rank: number;
  searches: number;
  items: SearchItem[];
}

export interface LocalSearchResult {
  items: SearchItem[];
  suggestYouTube: boolean;
}

export interface YouTubeSearchResult {
  items: SearchItem[];
  remaining: number;
  cached: boolean;
}

export type YouTubeSearchError =
  "quota_exhausted" | "user_limit" | "search_unavailable" | "youtube_error" | "unauthorized";
