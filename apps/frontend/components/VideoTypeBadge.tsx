import { videoTypeDotColor } from "@/lib/videoType";

export function VideoTypeBadge({ videoType }: Readonly<{ videoType: string }>) {
  return (
    <span className="badge">
      <span className="badge-dot" style={{ background: videoTypeDotColor(videoType) }} />
      {videoType}
    </span>
  );
}
