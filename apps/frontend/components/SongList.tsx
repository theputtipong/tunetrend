import { fetchSongs } from "@/lib/api";
import type { TabKey } from "@/lib/tabs";
import type { CountryCode } from "@/lib/countries";
import { countryLabel } from "@/lib/countries";
import { dictionaries, type Lang } from "@/lib/i18n";
import { SongRow } from "./SongRow";
import { StateMessage } from "./StateMessage";

export async function SongList({
  country,
  tab,
  category,
  lang,
}: Readonly<{ country: CountryCode; tab: TabKey; category?: string; lang: Lang }>) {
  const songs = await fetchSongs(country, tab, category);
  const t = dictionaries[lang];
  const countryName = countryLabel(country, lang);

  if (songs.length === 0) {
    return (
      <StateMessage
        variant="empty"
        title={t.songList.emptyTitle}
        description={t.songList.emptyDescription(countryName)}
      />
    );
  }

  return (
    <div className="flex flex-col gap-0.5 px-4 pb-8 pt-3 md:px-8">
      <p className="list-caption">{t.songList.caption(songs.length, countryName)}</p>
      {songs.map((song, index) => (
        <SongRow key={`${song.id}-${song.countryCode}`} song={song} rank={index + 1} lang={lang} />
      ))}
    </div>
  );
}
