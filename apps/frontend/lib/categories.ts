import type { Lang } from "./i18n";

export function resolveCategory(raw: string | undefined | null, available: ReadonlyArray<{ id: string }>): string {
  if (!raw) return "";
  return available.some((c) => c.id === raw) ? raw : "";
}

// รหัส id ต้องตรงกับ CategoryVideoConfigs ใน apps/backend/internal/domain/category_video.go
// หมวดที่ backend เพิ่มมาใหม่แต่ยังไม่มีคำแปลที่นี่ จะ fallback ไปใช้ title ภาษาอังกฤษแทน
const CATEGORY_LABELS_TH: Record<string, string> = {
  "1": "ภาพยนตร์และแอนิเมชัน",
  "2": "รถยนต์และยานพาหนะ",
  "20": "เกม",
  "24": "บันเทิง",
  "25": "ข่าวและการเมือง",
};

export function categoryLabel(id: string, title: string, lang: Lang): string {
  if (lang !== "th") return title;
  return CATEGORY_LABELS_TH[id] ?? title;
}
