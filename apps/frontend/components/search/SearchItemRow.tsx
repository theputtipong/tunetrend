"use client";

import Image from "next/image";
import Link from "next/link";
import type { ReactNode } from "react";
import { toPlaylistVideo } from "@/lib/playlist";
import type { SearchItem } from "@/types/search";
import { PlayIcon } from "../icons";
import { AddToPlaylistButton } from "../playlist/AddToPlaylistButton";

export function sendSearchHit(videoId: string, country: string, q: string) {
  try {
    void fetch("/api/search/hit", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ videoId, country, q }),
      keepalive: true,
    }).catch(() => {});
  } catch {}
}

export function SearchItemRow({
  item,
  href,
  rank,
  badges,
  statPrimary,
  statSecondary,
  hit,
}: Readonly<{
  item: SearchItem;
  href: string;
  rank?: number;
  badges?: ReactNode;
  statPrimary: ReactNode;
  statSecondary?: ReactNode;
  hit?: { country: string; q: string };
}>) {
  const isTop3 = rank !== undefined && rank <= 3;

  return (
    <div className="row w-full">
      <Link
        href={href}
        className="row-link text-left"
        onClick={hit ? () => sendSearchHit(item.id, hit.country, hit.q) : undefined}
      >
        {rank !== undefined ? (
          <span
            className={`rank w-5 text-[15px] md:w-8 ${
              isTop3 ? "rank--top md:text-[22px]" : "md:text-[18px]"
            }`}
          >
            {String(rank).padStart(2, "0")}
          </span>
        ) : null}

        <span className="thumb relative h-[47px] w-[84px] overflow-hidden md:h-[63px] md:w-[112px]">
          {item.thumbnailUrl ? (
            <Image
              src={item.thumbnailUrl}
              alt=""
              fill
              sizes="(min-width: 768px) 112px, 84px"
              className="object-cover"
            />
          ) : null}
          <span className="thumb-play absolute inset-0 flex items-center justify-center bg-black/0 transition-colors">
            <PlayIcon size={16} />
          </span>
        </span>

        <div className="min-w-0 flex-1">
          <div className="title">{item.title}</div>
          <div className="channel">{item.channelTitle}</div>
          {badges ? (
            <div className="mt-1.5 flex flex-wrap items-center gap-1.5">{badges}</div>
          ) : null}
          <div className="mt-1 text-[11.5px] text-[var(--text-3)] md:hidden">
            {statPrimary}
            {statSecondary ? <> · {statSecondary}</> : null}
          </div>
        </div>

        <div className="hidden w-[120px] flex-shrink-0 text-right md:block">
          <div className="views">{statPrimary}</div>
          {statSecondary ? <div className="time">{statSecondary}</div> : null}
        </div>
      </Link>
      <AddToPlaylistButton video={toPlaylistVideo(item)} />
    </div>
  );
}
