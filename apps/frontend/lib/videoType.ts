const VIDEO_TYPE_DOT_COLOR: Record<string, string> = {
  MV: "var(--accent)",
  Lyric: "oklch(74% 0.12 200)",
  "Audio Track": "oklch(74% 0.1 260)",
  Cover: "oklch(72% 0.12 320)",
  "Live Performance": "oklch(72% 0.14 140)",
  General: "oklch(60% 0.01 260)",
};

const FALLBACK_DOT_COLOR = "oklch(60% 0.01 260)";

export function videoTypeDotColor(videoType: string): string {
  return VIDEO_TYPE_DOT_COLOR[videoType] ?? FALLBACK_DOT_COLOR;
}
