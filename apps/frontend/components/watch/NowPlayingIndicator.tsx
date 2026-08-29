const BAR_COUNT = 10;

export function NowPlayingIndicator() {
  return (
    <div className="flex items-end gap-[1.6px]">
      {Array.from({ length: BAR_COUNT }).map((_, i) => (
        <span
          key={i}
          className="now-playing-bar block w-[2.2px] rounded-[2px] bg-[var(--accent)]"
          style={{ height: 14, animationDelay: `${(i / BAR_COUNT) * 0.9}s` }}
        />
      ))}
    </div>
  );
}
