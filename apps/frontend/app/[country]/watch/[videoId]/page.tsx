import type { CountryCode } from "@/lib/countries";
import { isValidCountry, countryLabel } from "@/lib/countries";
import { fetchDiscoverItems, fetchSongs } from "@/lib/api";
import { dictionaries } from "@/lib/i18n";
import { getLang } from "@/lib/i18n/server";
import { resolveTab } from "@/lib/tabs";
import { WatchPageClient } from "@/components/watch/WatchPageClient";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ country: string; videoId: string }>;
}) {
  const { country: rawCountry } = await params;
  const country = rawCountry.toUpperCase();
  const lang = await getLang();
  const label = isValidCountry(country) ? countryLabel(country, lang) : rawCountry;
  return { title: dictionaries[lang].countryPageTitle(label) };
}

export default async function WatchPage({
  params,
  searchParams,
}: Readonly<{
  params: Promise<{ country: string; videoId: string }>;
  searchParams: Promise<{ tab?: string }>;
}>) {
  const { country: rawCountry, videoId } = await params;
  const { tab: rawTab } = await searchParams;

  const country = rawCountry.toUpperCase() as CountryCode;
  const tab = resolveTab(rawTab);
  const lang = await getLang();

  // Related queue intentionally mirrors the mobile app: it always shows the
  // generic (non-category-filtered) list for the tab, even if the video being
  // watched came from a category-filtered view.
  const songs = await fetchSongs(country, tab);
  const current = songs.find((s) => s.id === videoId);

  // Discover carousel is a nice-to-have — never let it break the watch page.
  const discoverItems = await fetchDiscoverItems().catch(() => []);

  return (
    <WatchPageClient
      country={country}
      tab={tab}
      lang={lang}
      videoId={videoId}
      title={current?.title ?? ""}
      categoryId={current?.categoryId ?? ""}
      initialSongs={songs}
      discoverItems={discoverItems}
    />
  );
}
