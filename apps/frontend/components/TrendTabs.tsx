import Link from "next/link";
import { dictionaries, type Lang } from "@/lib/i18n";
import { PAGE_TABS, type PageTabKey } from "@/lib/tabs";

export function TrendTabs({
  country,
  active,
  category,
  lang,
}: Readonly<{ country: string; active: PageTabKey; category: string; lang: Lang }>) {
  const t = dictionaries[lang].tabs;

  const visibleTabs = category
    ? PAGE_TABS.filter((key) => key !== "mv" && key !== "searched")
    : PAGE_TABS;

  return (
    <div
      className="-mb-px flex min-w-0 gap-5 overflow-x-auto pb-px [scrollbar-width:none] md:gap-7"
      data-tour="tabs"
    >
      {visibleTabs.map((key) => {
        const href = category
          ? `/${country}?tab=${key}&category=${category}`
          : `/${country}?tab=${key}`;

        return (
          <Link
            key={key}
            href={href}
            className={`${key === active ? "tab tab--active" : "tab"} flex-shrink-0 whitespace-nowrap`}
          >
            {t[key]}
            {key === "searched" ? (
              <span className="tab-new">{dictionaries[lang].nav.newBadge}</span>
            ) : null}
          </Link>
        );
      })}
    </div>
  );
}
