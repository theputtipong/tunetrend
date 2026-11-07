import type { CountryCode } from "./countries";
import type { TabKey } from "./tabs";
import type { Song, TrendsResponse } from "@/types/song";
import type { Category, CategoriesResponse } from "@/types/category";
import type { DiscoverItem, DiscoverResponse } from "@/types/discover";
import type { LocalSearchResult, SearchedQuery } from "@/types/search";
import { forwardedClientHeaders } from "./session";

const API_BASE_URL = process.env.API_BASE_URL ?? "http://localhost:8080";

const ENDPOINT_BY_TAB: Record<TabKey, string> = {
  trending: "/trends",
  new: "/trends/new",
  mv: "/trends/mv",
};

const REVALIDATE_SECONDS = 3600;

export async function fetchSongs(
  country: CountryCode,
  tab: TabKey,
  category?: string,
): Promise<Song[]> {
  const params = new URLSearchParams({ country });
  if (category) params.set("category", category);

  const url = `${API_BASE_URL}${ENDPOINT_BY_TAB[tab]}?${params.toString()}`;

  const res = await fetch(url, { next: { revalidate: REVALIDATE_SECONDS } });

  if (res.status === 400 && category) {
    return fetchSongs(country, tab);
  }

  if (!res.ok) {
    throw new Error(`TuneTrend API responded with status ${res.status}`);
  }

  const body: TrendsResponse = await res.json();
  if (!body.success) {
    throw new Error(body.error ?? "TuneTrend API returned an error");
  }

  return body.data ?? [];
}

export async function fetchDiscoverItems(): Promise<DiscoverItem[]> {
  const url = `${API_BASE_URL}/discover`;

  const res = await fetch(url, { next: { revalidate: REVALIDATE_SECONDS } });
  if (!res.ok) {
    throw new Error(`TuneTrend API responded with status ${res.status}`);
  }

  const body: DiscoverResponse = await res.json();
  if (!body.success) {
    throw new Error(body.error ?? "TuneTrend API returned an error");
  }

  return body.data ?? [];
}

export async function fetchCategories(country: CountryCode): Promise<Category[]> {
  const url = `${API_BASE_URL}/categories?country=${country}`;

  const res = await fetch(url, { next: { revalidate: REVALIDATE_SECONDS } });
  if (!res.ok) {
    throw new Error(`TuneTrend API responded with status ${res.status}`);
  }

  const body: CategoriesResponse = await res.json();
  if (!body.success) {
    throw new Error(body.error ?? "TuneTrend API returned an error");
  }

  return body.data ?? [];
}

const SEARCHED_REVALIDATE_SECONDS = 300;
const OEMBED_REVALIDATE_SECONDS = 86400;

export async function fetchLocalSearch(
  q: string,
  country: CountryCode,
  category: string,
  clientIp?: string | null,
  track = true,
): Promise<LocalSearchResult> {
  const params = new URLSearchParams({ q, country, category });
  if (!track) params.set("track", "0");
  const res = await fetch(`${API_BASE_URL}/search?${params.toString()}`, {
    cache: "no-store",
    headers: forwardedClientHeaders(clientIp),
    signal: AbortSignal.timeout(10_000),
  });
  const body = await res.json().catch(() => null);
  if (!res.ok || !body?.success) {
    throw new Error(body?.error ?? `TuneTrend API responded with status ${res.status}`);
  }
  return { items: body.data?.items ?? [], suggestYouTube: Boolean(body.data?.suggestYouTube) };
}

export async function fetchSearchedQueries(country: CountryCode): Promise<SearchedQuery[]> {
  const res = await fetch(`${API_BASE_URL}/trends/searched-queries?country=${country}`, {
    next: { revalidate: SEARCHED_REVALIDATE_SECONDS },
  });
  const body = await res.json().catch(() => null);
  if (!res.ok || !body?.success) {
    throw new Error(body?.error ?? `TuneTrend API responded with status ${res.status}`);
  }
  return body.data ?? [];
}

export async function fetchVideoOEmbed(
  videoId: string,
): Promise<{ title: string; channelTitle: string } | null> {
  const url = `https://www.youtube.com/oembed?format=json&url=${encodeURIComponent(
    `https://www.youtube.com/watch?v=${videoId}`,
  )}`;
  try {
    const res = await fetch(url, {
      next: { revalidate: OEMBED_REVALIDATE_SECONDS },
      signal: AbortSignal.timeout(3000),
    });
    if (!res.ok) return null;
    const body = await res.json();
    if (typeof body?.title !== "string") return null;
    return {
      title: body.title,
      channelTitle: typeof body.author_name === "string" ? body.author_name : "",
    };
  } catch {
    return null;
  }
}
