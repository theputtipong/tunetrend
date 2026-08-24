"use client";

import Link from "next/link";
import { useEffect, useRef } from "react";
import { useSearchParams } from "next/navigation";
import { resolveCategory } from "@/lib/categories";
import { resolveTab } from "@/lib/tabs";
import type { Category } from "@/types/category";

const AUTO_ADVANCE_MS = 3000;

export function CategoryCarousel({
  country,
  categories,
}: Readonly<{ country: string; categories: Category[] }>) {
  const searchParams = useSearchParams();
  const tab = resolveTab(searchParams.get("tab"));
  const activeCategory = resolveCategory(searchParams.get("category"), categories);

  const trackRef = useRef<HTMLDivElement>(null);
  const pausedRef = useRef(false);

  useEffect(() => {
    const track = trackRef.current;
    if (!track) return;

    const id = setInterval(() => {
      if (pausedRef.current) return;

      const { scrollLeft, scrollWidth, clientWidth } = track;
      const atEnd = scrollLeft + clientWidth >= scrollWidth - 1;

      track.scrollTo({
        left: atEnd ? 0 : scrollLeft + clientWidth * 0.4,
        behavior: "smooth",
      });
    }, AUTO_ADVANCE_MS);

    return () => clearInterval(id);
  }, []);

  const options = [{ id: "", title: "เพลงฮิต" }, ...categories];

  return (
    <div
      ref={trackRef}
      className="flex items-center gap-2 overflow-x-auto px-4 pb-2 md:px-8"
      style={{ scrollSnapType: "x mandatory" }}
      onMouseEnter={() => {
        pausedRef.current = true;
      }}
      onMouseLeave={() => {
        pausedRef.current = false;
      }}
      onTouchStart={() => {
        pausedRef.current = true;
      }}
      onTouchEnd={() => {
        pausedRef.current = false;
      }}
      data-tour="category-carousel"
    >
      {options.map((cat) => {
        const active = cat.id === activeCategory;
        const href = cat.id
          ? `/${country}?tab=${tab}&category=${cat.id}`
          : `/${country}?tab=${tab}`;

        return (
          <Link
            key={cat.id || "all"}
            href={href}
            className={active ? "pill pill--active" : "pill"}
            style={{ scrollSnapAlign: "start" }}
          >
            {cat.title}
          </Link>
        );
      })}
    </div>
  );
}
