import { Suspense } from "react";
import Link from "next/link";
import { headers } from "next/headers";
import { fetchCategories, fetchLocalSearch } from "@/lib/api";
import { categoryLabel } from "@/lib/categories";
import { countryLabel, type CountryCode } from "@/lib/countries";
import { formatRelativeTime, formatViewCount } from "@/lib/format";
import { dictionaries, type Lang } from "@/lib/i18n";
import { getLang } from "@/lib/i18n/server";
import { getClientIp } from "@/lib/ratelimit";
import { isValidSearchQuery } from "@/lib/search";
import { MUSIC_CATEGORY_ID } from "@/lib/tabs";
import { Header } from "@/components/Header";
import { StateMessage } from "@/components/StateMessage";
import { VideoTypeBadge } from "@/components/VideoTypeBadge";
import { SearchItemRow } from "@/components/search/SearchItemRow";
import { YouTubeSection } from "@/components/search/YouTubeSection";
import type { SearchItem } from "@/types/search";

type SearchParams = Promise<{
  q?: string | string[];
  category?: string | string[];
  track?: string | string[];
}>;

function firstParam(value: string | string[] | undefined): string {
  return (Array.isArray(value) ? value[0] : value) ?? "";
}

export async function generateMetadata({ searchParams }: { searchParams: SearchParams }) {
  const q = firstParam((await searchParams).q).trim();
  const lang = await getLang();
  const t = dictionaries[lang].search;
  return { title: q ? t.metaTitle(q.slice(0, 100)) : t.promptTitle, robots: { index: false } };
}

export default async function SearchPage({
  params,
  searchParams,
}: Readonly<{ params: Promise<{ country: string }>; searchParams: SearchParams }>) {
  const { country: rawCountry } = await params;
  const raw = await searchParams;
  const country = rawCountry.toUpperCase() as CountryCode;
  const lang = await getLang();
  const t = dictionaries[lang];

  const q = firstParam(raw.q).trim();
  const categories = await fetchCategories(country).catch(() => []);
  const options = [
    { id: MUSIC_CATEGORY_ID, title: t.nav.musicCategory },
    ...categories.filter((c) => c.id !== MUSIC_CATEGORY_ID),
  ];
  const rawCategory = firstParam(raw.category);
  const category = options.some((o) => o.id === rawCategory) ? rawCategory : MUSIC_CATEGORY_ID;
  const lower = country.toLowerCase();
  const valid = isValidSearchQuery(q);

  return (
    <div className="mx-auto max-w-6xl">
      <Header
        activeCountry={country}
        countrySuffix={q ? `/search?${new URLSearchParams({ q }).toString()}` : ""}
        lang={lang}
        searchQuery={q}
        searchCategory={category === MUSIC_CATEGORY_ID ? undefined : category}
      />

      <nav
        aria-label={t.nav.searchLabel}
        className="flex flex-wrap items-center gap-2 px-4 py-3 md:px-8"
      >
        {options.map((opt) => {
          const qs = new URLSearchParams({ q });
          if (opt.id !== MUSIC_CATEGORY_ID) qs.set("category", opt.id);
          return (
            <Link
              key={opt.id}
              href={`/${lower}/search?${qs.toString()}`}
              className={opt.id === category ? "pill pill--active" : "pill"}
              aria-current={opt.id === category ? "page" : undefined}
            >
              {opt.id === MUSIC_CATEGORY_ID ? "🎵 " : ""}
              {categoryLabel(opt.id, opt.title, lang)}
            </Link>
          );
        })}
      </nav>

      <div className="border-t border-[var(--border)]" />

      {valid ? (
        <Suspense key={`${q}|${category}`} fallback={<ResultsSkeleton />}>
          <SearchResults
            q={q}
            country={country}
            category={category}
            lang={lang}
            track={firstParam(raw.track) !== "0"}
          />
        </Suspense>
      ) : (
        <StateMessage
          variant="empty"
          title={t.search.promptTitle}
          description={t.search.promptDescription}
        />
      )}
    </div>
  );
}

function chartBadge(item: SearchItem, lang: Lang) {
  const t = dictionaries[lang];
  if (item.source === "chart" && item.chartRank && item.chartTab) {
    return (
      <span className="badge badge--accent">
        {t.search.chartRankBadge(item.chartRank, t.tabs[item.chartTab])}
      </span>
    );
  }
  if (item.source === "searched") {
    return <span className="badge">{t.search.searchedBefore}</span>;
  }
  return null;
}

async function SearchResults({
  q,
  country,
  category,
  lang,
  track,
}: Readonly<{ q: string; country: CountryCode; category: string; lang: Lang; track: boolean }>) {
  const t = dictionaries[lang].search;
  let result;
  try {
    const clientIp = getClientIp(new Request("http://local", { headers: await headers() }));
    result = await fetchLocalSearch(q, country, category, clientIp, track);
  } catch {
    return <StateMessage variant="error" title={t.errorTitle} description={t.errorDescription} />;
  }

  const lower = country.toLowerCase();

  return (
    <div className="flex flex-col gap-0.5 px-4 pb-10 pt-1 md:px-8">
      <div className="sec-title">
        <b>{t.chartSection}</b>
        <span>{t.chartCount(result.items.length, countryLabel(country, lang))}</span>
      </div>

      {result.items.length === 0 ? (
        <p className="px-[14px] pb-2 text-sm text-[var(--text-3)]">
          <span className="font-semibold text-[var(--text-2)]">{t.noLocalTitle}</span> ·{" "}
          {t.noLocalDescription}
        </p>
      ) : (
        result.items.map((item) => (
          <SearchItemRow
            key={item.id}
            item={item}
            href={`/${lower}/watch/${item.id}`}
            hit={{ country, q }}
            badges={
              <>
                {chartBadge(item, lang)}
                {item.videoType ? <VideoTypeBadge videoType={item.videoType} /> : null}
              </>
            }
            statPrimary={formatViewCount(item.viewCount, lang)}
            statSecondary={formatRelativeTime(item.publishedAt, lang)}
          />
        ))
      )}

      <YouTubeSection
        q={q}
        country={country}
        category={category}
        suggestYouTube={result.suggestYouTube}
        excludeIds={result.items.map((i) => i.id)}
      />
    </div>
  );
}

function ResultsSkeleton() {
  return (
    <div className="flex flex-col gap-0.5 px-4 pt-3 md:px-8">
      {[0, 1, 2, 3].map((i) => (
        <div key={i} className="row animate-pulse">
          <span className="thumb h-[47px] w-[84px] md:h-[63px] md:w-[112px]" />
          <div className="min-w-0 flex-1">
            <div className="h-4 w-2/5 rounded bg-[var(--bg-elev-2)]" />
            <div className="mt-2 h-3 w-1/4 rounded bg-[var(--bg-elev-2)]" />
          </div>
        </div>
      ))}
    </div>
  );
}
