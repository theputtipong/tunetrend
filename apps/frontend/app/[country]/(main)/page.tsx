import type { CountryCode } from "@/lib/countries";
import { isValidCountry, countryLabel } from "@/lib/countries";
import { dictionaries } from "@/lib/i18n";
import { getLang } from "@/lib/i18n/server";
import { resolvePageTab } from "@/lib/tabs";
import { SongList } from "@/components/SongList";
import { SearchedList } from "@/components/SearchedList";

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
  const tab = resolvePageTab(rawTab);
  const lang = await getLang();

  if (tab === "searched") {
    return <SearchedList country={country} lang={lang} />;
  }

  const category = tab === "mv" ? "" : (rawCategory ?? "");

  return <SongList country={country} tab={tab} category={category} lang={lang} />;
}
