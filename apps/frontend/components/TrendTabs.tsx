import Link from "next/link";
import { dictionaries, type Lang } from "@/lib/i18n";
import { TABS, type TabKey } from "@/lib/tabs";

export function TrendTabs({
  country,
  active,
  lang,
}: Readonly<{ country: string; active: TabKey; lang: Lang }>) {
  const t = dictionaries[lang].tabs;

  return (
    <div className="flex gap-7" data-tour="tabs">
      {TABS.map((key) => (
        <Link
          key={key}
          href={`/${country}?tab=${key}`}
          className={key === active ? "tab tab--active" : "tab"}
        >
          {t[key]}
        </Link>
      ))}
    </div>
  );
}
