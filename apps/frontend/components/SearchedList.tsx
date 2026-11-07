import Link from "next/link";
import { fetchSearchedQueries } from "@/lib/api";
import type { CountryCode } from "@/lib/countries";
import { formatCompactNumber, formatRelativeTime, formatViewCount } from "@/lib/format";
import { dictionaries, type Lang } from "@/lib/i18n";
import { SearchItemRow } from "./search/SearchItemRow";
import { StateMessage } from "./StateMessage";
import { VideoTypeBadge } from "./VideoTypeBadge";

export async function SearchedList({
  country,
  lang,
}: Readonly<{ country: CountryCode; lang: Lang }>) {
  const groups = await fetchSearchedQueries(country);
  const t = dictionaries[lang].searched;
  const lower = country.toLowerCase();

  if (groups.length === 0) {
    return <StateMessage variant="empty" title={t.emptyTitle} description={t.emptyDescription} />;
  }

  return (
    <div className="flex flex-col gap-4 px-4 pb-8 pt-3 md:px-8">
      <p className="list-caption">{t.hint}</p>
      {groups.map((group) => {
        const searchHref = `/${lower}/search?${new URLSearchParams({ q: group.query, track: "0" }).toString()}`;
        return (
          <section key={group.query} className="query-group" aria-label={group.query}>
            <header className="query-group__head">
              <span className={`rank ${group.rank <= 3 ? "rank--top" : ""}`}>
                {String(group.rank).padStart(2, "0")}
              </span>
              <Link href={searchHref} className="query-group__term">
                {group.query}
              </Link>
              <span className="query-group__count">
                {t.searches(formatCompactNumber(group.searches))}
              </span>
              <Link href={searchHref} className="query-group__all">
                {t.viewAll} →
              </Link>
            </header>
            <div className="flex flex-col gap-0.5">
              {group.items.map((item) => (
                <SearchItemRow
                  key={item.id}
                  item={item}
                  href={`/${lower}/watch/${item.id}`}
                  badges={
                    item.videoType ? <VideoTypeBadge videoType={item.videoType} /> : undefined
                  }
                  statPrimary={formatViewCount(item.viewCount, lang)}
                  statSecondary={formatRelativeTime(item.publishedAt, lang)}
                />
              ))}
            </div>
          </section>
        );
      })}
    </div>
  );
}
