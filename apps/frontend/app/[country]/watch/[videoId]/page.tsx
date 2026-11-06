import type { CountryCode } from "@/lib/countries";
import { isValidCountry, countryLabel } from "@/lib/countries";
import { fetchDiscoverItems, fetchSongs, fetchVideoOEmbed } from "@/lib/api";
import { dictionaries } from "@/lib/i18n";
import { getLang } from "@/lib/i18n/server";
import { UUID_REGEX, VIDEO_ID_REGEX, playlistThumbnail } from "@/lib/playlist";
import { getMyPlaylist } from "@/lib/session";
import { resolveTab } from "@/lib/tabs";
import { WatchPageClient } from "@/components/watch/WatchPageClient";
import type { Song } from "@/types/song";

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

function shuffled<T>(items: T[]): T[] {
  const out = [...items];
  for (let i = out.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [out[i], out[j]] = [out[j], out[i]];
  }
  return out;
}

async function loadPlaylistQueue(
  playlistId: string | undefined,
  videoId: string,
  country: CountryCode,
  shuffle: boolean,
): Promise<{ id: string; name: string; songs: Song[] } | null> {
  if (!playlistId || !UUID_REGEX.test(playlistId)) return null;
  const playlist = await getMyPlaylist(playlistId);
  if (!playlist?.items?.length) return null;

  let songs: Song[] = playlist.items.map((item) => ({
    id: item.videoId,
    title: item.title,
    channelTitle: item.channelTitle,
    thumbnailUrl: playlistThumbnail(item.videoId, item.thumbnailUrl),
    viewCount: "",
    countryCode: country,
    categoryId: "",
    publishedAt: "",
    videoType: "",
  }));
  if (shuffle) {
    const rest = shuffled(songs.filter((s) => s.id !== videoId));
    const current = songs.find((s) => s.id === videoId);
    songs = current ? [current, ...rest] : rest;
  }
  return { id: playlist.id, name: playlist.name, songs };
}

export default async function WatchPage({
  params,
  searchParams,
}: Readonly<{
  params: Promise<{ country: string; videoId: string }>;
  searchParams: Promise<{ tab?: string; playlist?: string; shuffle?: string }>;
}>) {
  const { country: rawCountry, videoId } = await params;
  const { tab: rawTab, playlist: rawPlaylist, shuffle } = await searchParams;

  const country = rawCountry.toUpperCase() as CountryCode;
  const tab = resolveTab(rawTab);
  const lang = await getLang();

  const playlist = await loadPlaylistQueue(rawPlaylist, videoId, country, shuffle === "1");

  // The related queue is always the unfiltered list for the tab, even when the video came from a category view.
  const songs = playlist?.songs ?? (await fetchSongs(country, tab));
  const current = songs.find((s) => s.id === videoId);
  const fallback =
    !current && VIDEO_ID_REGEX.test(videoId) ? await fetchVideoOEmbed(videoId) : null;

  const discoverItems = await fetchDiscoverItems().catch(() => []);

  return (
    <WatchPageClient
      country={country}
      tab={tab}
      lang={lang}
      videoId={videoId}
      title={current?.title ?? fallback?.title ?? ""}
      channelTitle={current?.channelTitle ?? fallback?.channelTitle ?? ""}
      categoryId={current?.categoryId ?? ""}
      initialSongs={songs}
      discoverItems={discoverItems}
      playlist={playlist ? { id: playlist.id, name: playlist.name } : undefined}
    />
  );
}
