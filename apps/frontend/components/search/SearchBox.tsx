"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { dictionaries, type Lang } from "@/lib/i18n";

function isTypingTarget(target: EventTarget | null): boolean {
  if (!(target instanceof HTMLElement)) return false;
  const tag = target.tagName;
  return tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT" || target.isContentEditable;
}

export function SearchBox({
  country,
  lang,
  initialQuery = "",
  category,
}: Readonly<{ country: string; lang: Lang; initialQuery?: string; category?: string }>) {
  const router = useRouter();
  const inputRef = useRef<HTMLInputElement>(null);
  const [value, setValue] = useState(initialQuery);
  const [focused, setFocused] = useState(false);
  const t = dictionaries[lang].nav;

  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if (e.key !== "/" || e.metaKey || e.ctrlKey || e.altKey) return;
      if (isTypingTarget(e.target)) return;
      e.preventDefault();
      inputRef.current?.focus();
      inputRef.current?.select();
    }
    document.addEventListener("keydown", onKeyDown);
    return () => document.removeEventListener("keydown", onKeyDown);
  }, []);

  function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    const q = value.trim();
    if (!q) return;
    const params = new URLSearchParams({ q });
    if (category) params.set("category", category);
    router.push(`/${country.toLowerCase()}/search?${params.toString()}`);
  }

  return (
    <form role="search" onSubmit={onSubmit} className="search-box w-full">
      <svg
        width="15"
        height="15"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        aria-hidden="true"
      >
        <circle cx="11" cy="11" r="7" />
        <path d="m20 20-3.5-3.5" />
      </svg>
      <input
        ref={inputRef}
        type="search"
        name="q"
        value={value}
        maxLength={100}
        enterKeyHint="search"
        autoComplete="off"
        aria-label={t.searchLabel}
        placeholder={t.searchPlaceholder}
        onChange={(e) => setValue(e.target.value)}
        onFocus={() => setFocused(true)}
        onBlur={() => setFocused(false)}
        onKeyDown={(e) => {
          if (e.key === "Escape") {
            if (value) setValue("");
            else inputRef.current?.blur();
          }
        }}
      />
      <span className="kbd hidden sm:inline-flex" aria-hidden="true">
        {focused ? "esc" : "/"}
      </span>
    </form>
  );
}
