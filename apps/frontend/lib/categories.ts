import type { Lang } from "./i18n";

export function resolveCategory(
  raw: string | undefined | null,
  available: ReadonlyArray<{ id: string }>,
): string {
  if (!raw) return "";
  return available.some((c) => c.id === raw) ? raw : "";
}

const CATEGORY_LABELS_TH: Record<string, string> = {
  "1": "ภาพยนตร์และแอนิเมชัน",
  "2": "รถยนต์และยานพาหนะ",
  "15": "สัตว์เลี้ยงและสัตว์",
  "17": "กีฬา",
  "19": "การเดินทางและกิจกรรม",
  "20": "เกม",
  "22": "ผู้คนและบล็อก",
  "23": "ตลก",
  "24": "บันเทิง",
  "25": "ข่าวและการเมือง",
  "26": "วิธีทำและสไตล์",
  "27": "การศึกษา",
  "28": "วิทยาศาสตร์และเทคโนโลยี",
};

export function categoryLabel(id: string, title: string, lang: Lang): string {
  if (lang !== "th") return title;
  return CATEGORY_LABELS_TH[id] ?? title;
}
