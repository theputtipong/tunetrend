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

  return (
    <div className="flex gap-7" data-tour="tabs">
      {TABS.map((key) => {
        // "mv" ไม่รองรับการกรองตามหมวดหมู่ (ตาราง CategoryVideo ไม่มี videoType) จึงทิ้ง category ทิ้งเมื่อสลับไปแท็บนี้
        const href =
          key !== "mv" && category
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
