"use client";

import { useEffect, useRef, useState } from "react";
import Link from "next/link";
import { dictionaries, type Lang } from "@/lib/i18n";
import { InfoIcon, MoreIcon } from "./icons";
import { LanguageToggle } from "./LanguageToggle";
import { ReplayTourButton } from "./ReplayTourButton";
import { ThemeToggle } from "./ThemeToggle";

export function MobileMenu({
  lang,
  showAbout = true,
  showReplayTour = true,
}: Readonly<{ lang: Lang; showAbout?: boolean; showReplayTour?: boolean }>) {
  const [isOpen, setIsOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  const t = dictionaries[lang].nav;

  useEffect(() => {
    if (!isOpen) return;
    function onClickOutside(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setIsOpen(false);
    }
    document.addEventListener("mousedown", onClickOutside);
    return () => document.removeEventListener("mousedown", onClickOutside);
  }, [isOpen]);

  function close() {
    setIsOpen(false);
  }

  return (
    <div ref={ref} className="relative sm:hidden" data-tour="mobile-menu">
      <button
        type="button"
        onClick={() => setIsOpen((v) => !v)}
        className="theme-toggle"
        aria-label={t.menu}
        title={t.menu}
        aria-expanded={isOpen}
      >
        <MoreIcon />
      </button>

      {isOpen && (
        <div className="absolute right-0 top-full z-50 mt-2 min-w-[170px] rounded-xl border border-[var(--border)] bg-[var(--bg-elev)] p-1.5 shadow-lg">
          {showReplayTour && <ReplayTourButton lang={lang} variant="menu-item" onAction={close} />}
          <LanguageToggle variant="menu-item" onAction={close} />
          <ThemeToggle lang={lang} variant="menu-item" onAction={close} />
          {showAbout && (
            <Link href="/about" onClick={close} className="menu-item">
              <InfoIcon size={16} />
              {t.about}
            </Link>
          )}
        </div>
      )}
    </div>
  );
}
