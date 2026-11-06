import { fetchSearchedTrends } from "@/lib/api";
import type { CountryCode } from "@/lib/countries";
import { formatCompactNumber } from "@/lib/format";
import { dictionaries, type Lang } from "@/lib/i18n";
import { SearchItemRow } from "./search/SearchItemRow";
import { StateMessage } from "./StateMessage";
import { VideoTypeBadge } from "./VideoTypeBadge";

export async function SearchedList({
  country,
  lang,
}: Readonly<{ country: CountryCode; lang: Lang }>) {
  const trends = await fetchSearchedTrends(country);
  const t = dictionaries[lang].searched;

  if (trends.length === 0) {
    return <StateMessage variant="empty" title={t.emptyTitle} description={t.emptyDescription} />;
  }

  return (
    <div className="flex flex-col gap-0.5 px-4 pb-8 pt-3 md:px-8">
      <p className="list-caption">{t.hint}</p>
      {trends.map((trend) => {
        const change = trend.isNew ? (
          <span className="rank-up">{t.isNew}</span>
        ) : trend.rankChange > 0 ? (
          <span className="rank-up">{t.rankUp(trend.rankChange)}</span>
        ) : trend.rankChange < 0 ? (
          <span className="rank-down">{t.rankDown(-trend.rankChange)}</span>
        ) : undefined;

        return (
          <SearchItemRow
            key={trend.id}
            item={trend}
            rank={trend.rank}
            href={`/${country.toLowerCase()}/watch/${trend.id}`}
            badges={trend.videoType ? <VideoTypeBadge videoType={trend.videoType} /> : undefined}
            statPrimary={t.searches(formatCompactNumber(trend.hits))}
            statSecondary={change}
          />
        );
      })}
    </div>
  );
}
