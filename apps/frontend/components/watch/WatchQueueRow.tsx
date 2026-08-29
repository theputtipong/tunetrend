import { forwardRef } from "react";
import Image from "next/image";
import type { Song } from "@/types/song";
import { formatRelativeTime, formatViewCount } from "@/lib/format";
import type { Lang } from "@/lib/i18n";
import { PlayIcon } from "../icons";
import { VideoTypeBadge } from "../VideoTypeBadge";
import { NowPlayingIndicator } from "./NowPlayingIndicator";

export const WatchQueueRow = forwardRef<
  HTMLButtonElement,
  Readonly<{
    song: Song;
    rank: number;
    lang: Lang;
    isPlaying: boolean;
    onSelect: () => void;
  }>
>(function WatchQueueRow({ song, rank, lang, isPlaying, onSelect }, ref) {
  const views = formatViewCount(song.viewCount, lang);
  const time = formatRelativeTime(song.publishedAt, lang);

  return (
    <button
      ref={ref}
      type="button"
      onClick={onSelect}
      disabled={isPlaying}
      className="row relative w-full appearance-none border-0 bg-transparent text-left"
    >
      <span className="rank w-5 flex-shrink-0 text-[15px]">{String(rank).padStart(2, "0")}</span>

      <span className="thumb relative h-[47px] w-[84px] flex-shrink-0 overflow-hidden md:h-[63px] md:w-[112px]">
        <Image src={song.thumbnailUrl} alt="" fill sizes="112px" className="object-cover" />
        {!isPlaying && (
          <span className="thumb-play absolute inset-0 flex items-center justify-center bg-black/0 transition-colors">
            <PlayIcon size={16} />
          </span>
        )}
      </span>

      <div className="min-w-0 flex-1">
        <div className="title">{song.title}</div>
        <div className="channel">{song.channelTitle}</div>
        <div className="mt-1.5 flex flex-wrap items-center gap-2">
          <VideoTypeBadge videoType={song.videoType} />
          <span className="text-[11.5px] text-[var(--text-3)]">
            {views} · {time}
          </span>
        </div>
      </div>

      {isPlaying && (
        <span
          className="absolute inset-0 flex items-center justify-center rounded-[inherit]"
          style={{ background: "var(--now-playing-overlay)" }}
        >
          <NowPlayingIndicator />
        </span>
      )}
    </button>
  );
});
