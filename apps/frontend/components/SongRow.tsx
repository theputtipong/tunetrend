import Link from "next/link";
import Image from "next/image";
import type { Song } from "@/types/song";
import { formatRelativeTime, formatViewCount } from "@/lib/format";
import type { Lang } from "@/lib/i18n";
import type { TabKey } from "@/lib/tabs";
import { PlayIcon } from "./icons";
import { VideoTypeBadge } from "./VideoTypeBadge";

export function SongRow({
  song,
  rank,
  lang,
  country,
  tab,
  category,
}: Readonly<{
  song: Song;
  rank: number;
  lang: Lang;
  country: string;
  tab: TabKey;
  category?: string;
}>) {
  const isTop3 = rank <= 3;
  const views = formatViewCount(song.viewCount, lang);
  const time = formatRelativeTime(song.publishedAt, lang);

  const params = new URLSearchParams({ tab });
  if (category) params.set("category", category);
  const href = `/${country}/watch/${song.id}?${params.toString()}`;

  return (
    <Link href={href} className="row w-full appearance-none border-0 bg-transparent text-left">
      <span
        className={`rank w-5 text-[15px] md:w-8 ${
          isTop3 ? "rank--top md:text-[22px]" : "md:text-[18px]"
        }`}
      >
        {String(rank).padStart(2, "0")}
      </span>

      <span className="thumb relative h-[47px] w-[84px] overflow-hidden md:h-[63px] md:w-[112px]">
        <Image src={song.thumbnailUrl} alt="" fill sizes="112px" className="object-cover" />
        <span className="thumb-play absolute inset-0 flex items-center justify-center bg-black/0 transition-colors">
          <PlayIcon size={16} />
        </span>
      </span>

      <div className="min-w-0 flex-1">
        <div className="title">{song.title}</div>
        <div className="channel">{song.channelTitle}</div>
        <div className="mt-1.5 flex flex-wrap items-center gap-2 md:hidden">
          <VideoTypeBadge videoType={song.videoType} />
          <span className="text-[11.5px] text-[var(--text-3)]">
            {views} · {time}
          </span>
        </div>
      </div>

      <span className="hidden flex-shrink-0 md:inline-flex">
        <VideoTypeBadge videoType={song.videoType} />
      </span>

      <div className="hidden w-[120px] flex-shrink-0 text-right md:block">
        <div className="views">{views}</div>
        <div className="time">{time}</div>
      </div>
    </Link>
  );
}
