export type TabKey = "trending" | "new" | "mv";

export type PageTabKey = TabKey | "searched";

export const TABS: TabKey[] = ["trending", "new", "mv"];

export const PAGE_TABS: PageTabKey[] = [...TABS, "searched"];

export const MUSIC_CATEGORY_ID = "10";

const VALID_TABS = new Set<TabKey>(TABS);
const VALID_PAGE_TABS = new Set<PageTabKey>(PAGE_TABS);

export function resolveTab(rawTab: string | null | undefined): TabKey {
  return VALID_TABS.has(rawTab as TabKey) ? (rawTab as TabKey) : "trending";
}

export function resolvePageTab(rawTab: string | null | undefined): PageTabKey {
  return VALID_PAGE_TABS.has(rawTab as PageTabKey) ? (rawTab as PageTabKey) : "trending";
}
