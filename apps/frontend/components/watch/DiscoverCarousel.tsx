"use client";

import { useEffect, useRef } from "react";
import Image from "next/image";
import Link from "next/link";
import { dictionaries, type Lang } from "@/lib/i18n";
import type { DiscoverItem } from "@/types/discover";

const AUTO_ADVANCE_MS = 3000;

export function DiscoverCarousel({ items, lang }: Readonly<{ items: DiscoverItem[]; lang: Lang }>) {
  const t = dictionaries[lang].watch;
  const scrollerRef = useRef<HTMLDivElement>(null);
  const pausedRef = useRef(false);

  useEffect(() => {
    const scroller = scrollerRef.current;
    if (!scroller || items.length === 0) return;

    const timer = setInterval(() => {
      if (pausedRef.current) return;

      const card = scroller.querySelector<HTMLElement>("[data-card]");
      const step = card ? card.offsetWidth + 12 : 220;
      const atEnd = scroller.scrollLeft + scroller.clientWidth >= scroller.scrollWidth - 4;

      scroller.scrollTo({
        left: atEnd ? 0 : scroller.scrollLeft + step,
        behavior: "smooth",
      });
    }, AUTO_ADVANCE_MS);

    return () => clearInterval(timer);
  }, [items.length]);

  if (items.length === 0) return null;

  return (
    <div className="w-full bg-[var(--bg)] px-4 py-4 md:px-6">
      <h2 className="mb-3 text-sm font-semibold text-[var(--text-2)]">{t.discoverHeading}</h2>
      <div
        ref={scrollerRef}
        onMouseEnter={() => (pausedRef.current = true)}
        onMouseLeave={() => (pausedRef.current = false)}
        onTouchStart={() => (pausedRef.current = true)}
        onTouchEnd={() => (pausedRef.current = false)}
        className="flex gap-3 overflow-x-auto scroll-smooth pb-1 [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
      >
        {items.map((item) => (
          <Link
            key={`${item.id}-${item.categoryId}`}
            data-card
            href={`/${item.countryCode.toLowerCase()}?tab=trending&category=${item.categoryId}`}
            className="relative h-32 w-56 flex-shrink-0 overflow-hidden rounded-xl bg-black"
          >
            <Image
              src={item.thumbnailUrl}
              alt=""
              fill
              sizes="224px"
              className="object-cover opacity-80"
            />
            <span className="absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/90 to-transparent px-2.5 pb-2 pt-6 text-[11px] font-medium text-white/90">
              {item.title}
            </span>
            <span className="absolute right-1.5 top-1.5 rounded-full bg-black/70 px-2 py-0.5 text-[10.5px] font-semibold text-white">
              {item.categoryLabel}
            </span>
          </Link>
        ))}
      </div>
    </div>
  );
}
