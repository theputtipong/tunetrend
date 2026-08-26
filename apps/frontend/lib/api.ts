import type { CountryCode } from "./countries";
import type { TabKey } from "./tabs";
import type { Song, TrendsResponse } from "@/types/song";
import type { Category, CategoriesResponse } from "@/types/category";

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
