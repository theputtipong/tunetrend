import Link from "next/link";
import { dictionaries, type Lang } from "@/lib/i18n";
import { TABS, type TabKey } from "@/lib/tabs";

export function TrendTabs({
  country,
  active,
  category,
  lang,
}: Readonly<{ country: string; active: TabKey; category: string; lang: Lang }>) {
  const t = dictionaries[lang].tabs;

  const visibleTabs = category ? TABS.filter((key) => key !== "mv") : TABS;

  return (
    <div className="flex gap-7" data-tour="tabs">
      {visibleTabs.map((key) => {
        const href = category
          ? `/${country}?tab=${key}&category=${category}`
          : `/${country}?tab=${key}`;

        return (
          <Link key={key} href={href} className={key === active ? "tab tab--active" : "tab"}>
            {t[key]}
          </Link>
        );
      })}
    </div>
  );
}
