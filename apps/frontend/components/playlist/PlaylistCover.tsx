import Image from "next/image";
import type { PlaylistItem } from "@/types/playlist";
import { LogoMark } from "../icons";

export function PlaylistCover({
  items,
  className = "",
}: Readonly<{ items?: PlaylistItem[]; className?: string }>) {
  const thumbs = (items ?? []).filter((i) => i.thumbnailUrl).slice(0, 4);

  if (thumbs.length < 4) {
    return (
      <div className={`pl-cover flex items-center justify-center ${className}`}>
        {thumbs[0] ? (
          <span className="relative h-full w-full">
            <Image
              src={thumbs[0].thumbnailUrl}
              alt=""
              fill
              sizes="200px"
              className="object-cover"
            />
          </span>
        ) : (
          <span className="opacity-90 [--accent:var(--accent-ink)] [--bg:transparent]">
            <LogoMark size={40} />
          </span>
        )}
      </div>
    );
  }

  return (
    <div className={`pl-cover mosaic ${className}`}>
      {thumbs.map((item) => (
        <span key={item.videoId} className="relative">
          <Image src={item.thumbnailUrl} alt="" fill sizes="100px" className="object-cover" />
        </span>
      ))}
    </div>
  );
}
