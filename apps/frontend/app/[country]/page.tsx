import type { CountryCode } from "@/lib/countries";
import { isValidCountry, countryLabel } from "@/lib/countries";
import { dictionaries } from "@/lib/i18n";
import { getLang } from "@/lib/i18n/server";
import { resolveTab } from "@/lib/tabs";
import { SongList } from "@/components/SongList";

export async function generateMetadata({ params }: { params: Promise<{ country: string }> }) {
  const { country: rawCountry } = await params;
  const country = rawCountry.toUpperCase();
  const lang = await getLang();
  const label = isValidCountry(country) ? countryLabel(country, lang) : rawCountry;
  return { title: dictionaries[lang].countryPageTitle(label) };
}

export default async function CountryTrendsPage({
  params,
  searchParams,
}: Readonly<{
  params: Promise<{ country: string }>;
  searchParams: Promise<{ tab?: string; category?: string }>;
}>) {
  const { country: rawCountry } = await params;
  const { tab: rawTab, category: rawCategory } = await searchParams;

  const country = rawCountry.toUpperCase() as CountryCode;
  const tab = resolveTab(rawTab);
  // แท็บ "mv" รองรับเฉพาะหมวดเพลง (ตาราง Song เท่านั้นที่มี videoType) ถ้ามี category
  // ติดมาจาก URL ที่ไม่ถูกต้อง (เช่น bookmark เก่า) ให้ทิ้งไปแทนที่จะส่งไปหา backend เปล่าๆ
  // ไม่ validate กับ /categories ที่นี่ (ปล่อยให้ backend ตัดสินแล้ว fetchSongs fallback เอง) —
  // เพื่อให้ fetch เพลง/วิดีโอที่นี่ยิงพร้อมกับ fetchCategories ของ NavBar ใน layout.tsx ได้เลย
  // ไม่ต้องรอ /categories เสร็จก่อน ลด TTFB ลง 1 รอบ round-trip
  const category = tab === "mv" ? "" : (rawCategory ?? "");
  const lang = await getLang();

  return <SongList country={country} tab={tab} category={category} lang={lang} />;
}
