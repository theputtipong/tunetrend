import type { PlaylistVideo } from "@/types/playlist";

export const VIDEO_ID_REGEX = /^[A-Za-z0-9_-]{11}$/;
export const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

export const MAX_PLAYLISTS = 50;
export const MAX_PLAYLIST_ITEMS = 200;
export const MAX_PLAYLIST_NAME_LEN = 60;
export const MAX_ITEM_TITLE_LEN = 200;
export const MAX_ITEM_CHANNEL_LEN = 100;

const THUMBNAIL_PREFIX = "https://i.ytimg.com/";

export function runeLength(value: string): number {
  return Array.from(value).length;
}

export function isValidPlaylistName(raw: unknown): raw is string {
  if (typeof raw !== "string") return false;
  const len = runeLength(raw.trim());
  return len >= 1 && len <= MAX_PLAYLIST_NAME_LEN;
}

export function truncateRunes(value: string, max: number): string {
  const chars = Array.from(value.trim());
  return chars.length > max ? chars.slice(0, max).join("") : chars.join("");
}

export function playlistThumbnail(videoId: string, thumbnailUrl?: string): string {
  if (thumbnailUrl && thumbnailUrl.startsWith(THUMBNAIL_PREFIX) && thumbnailUrl.length <= 300) {
    return thumbnailUrl;
  }
  return `${THUMBNAIL_PREFIX}vi/${videoId}/mqdefault.jpg`;
}

export function toPlaylistVideo(video: {
  id: string;
  title: string;
  channelTitle: string;
  thumbnailUrl?: string;
}): PlaylistVideo {
  return {
    videoId: video.id,
    title: truncateRunes(video.title, MAX_ITEM_TITLE_LEN),
    channelTitle: truncateRunes(video.channelTitle, MAX_ITEM_CHANNEL_LEN),
    thumbnailUrl: playlistThumbnail(video.id, video.thumbnailUrl),
  };
}

const UNAVAILABLE_KEY = "tunetrend-unavailable-videos";
const UNAVAILABLE_EVENT = "tunetrend-unavailable-change";
const MAX_UNAVAILABLE = 500;

export function readUnavailableVideos(): string {
  try {
    return localStorage.getItem(UNAVAILABLE_KEY) ?? "";
  } catch {
    return "";
  }
}

export function markVideoUnavailable(videoId: string) {
  try {
    const ids = readUnavailableVideos().split(",").filter(Boolean);
    if (ids.includes(videoId)) return;
    ids.push(videoId);
    localStorage.setItem(UNAVAILABLE_KEY, ids.slice(-MAX_UNAVAILABLE).join(","));
    window.dispatchEvent(new Event(UNAVAILABLE_EVENT));
  } catch {}
}

export function subscribeUnavailable(callback: () => void) {
  window.addEventListener(UNAVAILABLE_EVENT, callback);
  window.addEventListener("storage", callback);
  return () => {
    window.removeEventListener(UNAVAILABLE_EVENT, callback);
    window.removeEventListener("storage", callback);
  };
}
