const SKELETON_ROWS = 8;

export default function CountryTrendsLoading() {
  return (
    <div className="flex flex-col gap-0.5 px-4 pb-8 pt-1 md:px-8">
      {Array.from({ length: SKELETON_ROWS }).map((_, i) => (
        <div key={i} className="row animate-pulse">
          <span className="h-[18px] w-5 rounded bg-[var(--bg-elev-2)] md:w-8" />
          <span className="thumb h-[47px] w-[84px] md:h-[63px] md:w-[112px]" />
          <div className="min-w-0 flex-1 space-y-2">
            <div className="h-[13px] w-2/5 rounded bg-[var(--bg-elev-2)]" />
            <div className="h-[11px] w-1/4 rounded bg-[var(--bg-elev-2)]" />
          </div>
        </div>
      ))}
    </div>
  );
}
