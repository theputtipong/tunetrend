import { dictionaries, type Lang } from "./i18n";

const compactViews = new Intl.NumberFormat("en-US", {
  notation: "compact",
  maximumFractionDigits: 1,
});

export function formatViewCount(viewCount: string, lang: Lang = "en"): string {
  const n = Number(viewCount);
  const suffix = dictionaries[lang].format.views;
  if (!Number.isFinite(n)) return `— ${suffix}`;
  return `${compactViews.format(n)} ${suffix}`;
}

export function formatRelativeTime(iso: string, lang: Lang = "en"): string {
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) return "";

  const t = dictionaries[lang].format;
  const diffMinutes = Math.round((Date.now() - date.getTime()) / 60_000);
  if (diffMinutes < 60) return t.minutesAgo(Math.max(diffMinutes, 0));

  const diffHours = Math.round(diffMinutes / 60);
  if (diffHours < 24) return t.hoursAgo(diffHours);

  const diffDays = Math.round(diffHours / 24);
  if (diffDays < 7) return t.daysAgo(diffDays);

  const diffWeeks = Math.round(diffDays / 7);
  return t.weeksAgo(diffWeeks);
}
