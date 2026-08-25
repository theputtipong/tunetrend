import Link from "next/link";
import { categoryLabel } from "@/lib/categories";
import type { TabKey } from "@/lib/tabs";
import { dictionaries, type Lang } from "@/lib/i18n";
import type { Category } from "@/types/category";

export function CategoryFilter({
  country,
  tab,
  active,
  categories,
  lang,
}: Readonly<{
  country: string;
  tab: TabKey;
  active: string;
  categories: Category[];
  lang: Lang;
}>) {
  const options = [{ id: "", title: dictionaries[lang].nav.musicCategory }, ...categories];

  return (
    <div
      // min-h กันไว้สำหรับ 2 บรรทัด (pill min-height 36px + gap 8px + padding 24px)
      // ป้องกัน layout shift ถ้า label ยาวขึ้น (เช่นภาษาไทย) แล้วดัน pill ตกบรรทัดใหม่หลัง first paint
      className="flex min-h-[104px] flex-wrap items-center gap-2 px-4 py-3 md:px-8"
      data-tour="category-filter"
    >
      {options.map((cat) => {
        const isActive = cat.id === active;
        // แท็บ "mv" มีเฉพาะหมวดเพลง ถ้ากำลังอยู่แท็บ mv แล้วเลือกหมวดอื่น ต้องสลับกลับไปแท็บ
        // trending ให้ด้วย ไม่งั้นจะได้ URL ที่ผสมกันแบบไม่มีความหมาย (?tab=mv&category=20)
        const targetTab = cat.id && tab === "mv" ? "trending" : tab;
        const href = cat.id
          ? `/${country}?tab=${targetTab}&category=${cat.id}`
          : `/${country}?tab=${tab}`;

        return (
          <Link
            key={cat.id || "music"}
            href={href}
            className={isActive ? "pill pill--active" : "pill"}
          >
            {categoryLabel(cat.id, cat.title, lang)}
          </Link>
        );
      })}
    </div>
  );
}
