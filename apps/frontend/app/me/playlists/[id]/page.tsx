import { headers } from "next/headers";
import { detectCountryFromAcceptLanguage } from "@/lib/countries";
import { dictionaries } from "@/lib/i18n";
import { getLang } from "@/lib/i18n/server";
import { UUID_REGEX } from "@/lib/playlist";
import { getCurrentUser, getMyPlaylist } from "@/lib/session";
import { StateMessage } from "@/components/StateMessage";
import { LoginRequired } from "@/components/playlist/LoginRequired";
import { PlaylistDetailClient } from "@/components/playlist/PlaylistDetailClient";

export async function generateMetadata() {
  const lang = await getLang();
  return { title: dictionaries[lang].playlist.metaTitle, robots: { index: false } };
}

export default async function PlaylistDetailPage({
  params,
}: Readonly<{ params: Promise<{ id: string }> }>) {
  const { id } = await params;
  const user = await getCurrentUser();
  if (!user) return <LoginRequired />;

  const lang = await getLang();
  const playlist = UUID_REGEX.test(id) ? await getMyPlaylist(id) : null;
  if (!playlist) {
    const t = dictionaries[lang].playlist;
    return <StateMessage variant="empty" title={t.notFoundTitle} description={t.notFoundBody} />;
  }

  const country = detectCountryFromAcceptLanguage((await headers()).get("accept-language"));
  return <PlaylistDetailClient key={playlist.id} initialPlaylist={playlist} country={country} />;
}
