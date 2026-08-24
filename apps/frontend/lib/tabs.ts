export type TabKey = "trending" | "new" | "mv";

export const TABS: TabKey[] = ["trending", "new", "mv"];

const VALID_TABS = new Set<TabKey>(TABS);

export function resolveTab(rawTab: string | null | undefined): TabKey {
  return VALID_TABS.has(rawTab as TabKey) ? (rawTab as TabKey) : "trending";
}
