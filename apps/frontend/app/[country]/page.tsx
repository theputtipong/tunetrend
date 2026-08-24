import type { CountryCode } from "@/lib/countries";
import { isValidCountry, countryLabel } from "@/lib/countries";
import { dictionaries } from "@/lib/i18n";
import { getLang } from "@/lib/i18n/server";
import { resolveTab } from "@/lib/tabs";
import { resolveCategory } from "@/lib/categories";
import { fetchCategories } from "@/lib/api";
import { SongList } from "@/components/SongList";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ country: string }>;
}) {
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
  const categories = await fetchCategories(country);
  const category = resolveCategory(rawCategory, categories);
  const lang = await getLang();

  return <SongList country={country} tab={tab} category={category} lang={lang} />;
}
