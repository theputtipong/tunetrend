import type { ImageLoaderProps } from "next/image";

const YT_THUMB_PATTERN = /^https:\/\/(?:i\.ytimg\.com|img\.youtube\.com)\/vi\/([A-Za-z0-9_-]+)\//;

const SMALL_FILE = "mqdefault.jpg";
const SMALL_MAX_WIDTH = 320;
const LARGE_FILE = "hqdefault.jpg";

export default function youtubeImageLoader({ src, width }: ImageLoaderProps): string {
  const match = YT_THUMB_PATTERN.exec(src);
  if (!match) return src;

  const file = width <= SMALL_MAX_WIDTH ? SMALL_FILE : LARGE_FILE;
  return `https://i.ytimg.com/vi/${match[1]}/${file}`;
}
