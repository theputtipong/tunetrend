import { dictionaries } from "@/lib/i18n";
import { getLang } from "@/lib/i18n/server";
import { getCurrentUser, getMyPlaylists } from "@/lib/session";
import { StateMessage } from "@/components/StateMessage";
import { LoginRequired } from "@/components/playlist/LoginRequired";
import { PlaylistsClient } from "@/components/playlist/PlaylistsClient";

export async function generateMetadata() {
  const lang = await getLang();
  return { title: dictionaries[lang].playlist.metaTitle, robots: { index: false } };
}

export default async function MyPlaylistsPage() {
  const user = await getCurrentUser();
  if (!user) return <LoginRequired />;

  const lang = await getLang();
  const playlists = await getMyPlaylists();
  if (!playlists) {
    const t = dictionaries[lang].playlist;
    return <StateMessage variant="error" title={t.actionFailed} description={t.errorGeneric} />;
  }

  return <PlaylistsClient initialPlaylists={playlists} />;
}
