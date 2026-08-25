const SKELETON_ROWS = 8;

export default function CountryTrendsLoading() {
  return (
    <div className="flex flex-col gap-0.5 px-4 pb-8 pt-3 md:px-8">
      <div className="list-caption animate-pulse">
        <span className="inline-block h-[11px] w-40 rounded bg-[var(--bg-elev-2)]" />
      </div>

      {Array.from({ length: SKELETON_ROWS }).map((_, i) => (
        <div key={i} className="row animate-pulse">
          <span className="h-[18px] w-5 rounded bg-[var(--bg-elev-2)] md:h-[22px] md:w-8" />

          <span className="thumb h-[47px] w-[84px] md:h-[63px] md:w-[112px]" />

          {/* บาร์สูงตามที่วัดจริงจาก .title/.channel/.badge/.views/.time (line-height ≈ 1.5×font-size
              ของฟอนต์นี้ ไม่ใช่ font-size ตรงๆ) ไม่งั้น mobile row จะเตี้ยกว่าของจริง ~23.5px */}
          <div className="min-w-0 flex-1">
            <div className="h-[22.5px] w-2/5 rounded bg-[var(--bg-elev-2)]" />
            <div className="mt-[2px] h-[19.5px] w-1/4 rounded bg-[var(--bg-elev-2)]" />

            <div className="mt-1.5 flex flex-wrap items-center gap-2 md:hidden">
              <span className="h-[28.5px] w-20 rounded-full bg-[var(--bg-elev-2)]" />
              <span className="h-[10px] w-16 rounded bg-[var(--bg-elev-2)]" />
            </div>
          </div>

          <span className="hidden h-[28.5px] w-20 flex-shrink-0 rounded-full bg-[var(--bg-elev-2)] md:inline-flex" />

          <div className="hidden w-[120px] flex-shrink-0 text-right md:block">
            <div className="ml-auto h-[20.25px] w-16 rounded bg-[var(--bg-elev-2)]" />
            <div className="ml-auto mt-[2px] h-[18px] w-12 rounded bg-[var(--bg-elev-2)]" />
          </div>
        </div>
      ))}
    </div>
  );
}
